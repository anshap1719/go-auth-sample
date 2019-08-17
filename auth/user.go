package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/auth/models"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/email"
	"gigglesearch.org/giggle-auth/utils/secrets"
	"github.com/globalsign/mgo"
	"github.com/goadesign/goa"
	"github.com/simukti/emailcheck"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const emailValidateExpiration = 7 * 24 * time.Hour

var ErrInvalidRecaptcha = errors.New("Invalid recaptcha response")

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	auth.JWTSecurity
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, jwtSec auth.JWTSecurity) *UserController {
	return &UserController{
		Controller:  service.NewController("UserController"),
		JWTSecurity: jwtSec,
	}
}

func (c *UserController) AddPlugin(ctx *app.AddPluginUserContext) error {
	return nil
}

// Deactivate runs the deactivate action.
func (c *UserController) Deactivate(ctx *app.DeactivateUserContext) error {
	// UserController_Deactivate: start_implement

	var uID string

	if *ctx.Admin {
		uID = *ctx.ID
	} else {
		uID = c.GetUserID(ctx.Request)
	}

	user, err := database.GetUser(ctx, uID)
	if err == mgo.ErrNotFound {
		return ctx.InternalServerError(goa.ErrNotFound(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	type access struct {
		query    func(context.Context, string) (string, error)
		notFound error
		deletion func(context.Context, string) error
	}

	getAccountItems := []access{
		{
			query:    database.QueryGoogleAccountUser,
			notFound: database.ErrGoogleAccountNotFound,
			deletion: database.DeleteGoogleAccount,
		},
		{
			query:    database.QueryFacebookAccountUser,
			notFound: database.ErrFacebookAccountNotFound,
			deletion: database.DeleteFacebookAccount,
		},
		{
			query:    database.QueryTwitterAccountUser,
			notFound: database.ErrTwitterAccountNotFound,
			deletion: database.DeleteTwitterAccount,
		},
		{
			query:    database.QueryLinkedinAccountUser,
			notFound: database.ErrLinkedinAccountNotFound,
			deletion: database.DeleteLinkedinAccount,
		},
		{
			query:    database.QueryMicrosoftAccountUser,
			notFound: database.ErrMicrosoftAccountNotFound,
			deletion: database.DeleteMicrosoftAccount,
		},
		{
			query:    database.QueryPasswordLoginFromID,
			notFound: database.ErrPasswordLoginNotFound,
			deletion: database.DeletePasswordLogin,
		},
	}

	for _, v := range getAccountItems {
		k, err := v.query(ctx, uID)
		if err == nil {
			err = v.deletion(ctx, k)
			if err != nil {
				return ctx.InternalServerError(goa.ErrInternal(err))
			}
		} else if err != v.notFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	}

	// @TODO: Send Notification Via Email

	_ = struct {
		UserAbout string
		Type      string
	}{
		UserAbout: uID,
		Type:      "user-disabled",
	}

	toName := user.FirstName + " " + user.LastName
	toMail := user.Email
	textContent := "Your account has been disabled"
	htmlContent := "Your account has been disabled"
	subject := "Your account has been disabled"

	if err := email.SendMail(subject, toName, toMail, textContent, htmlContent); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))

	// UserController_Deactivate: end_implement
}

// GetAllUsers runs the get-all-users action.
func (c *UserController) GetAllUsers(ctx *app.GetAllUsersUserContext) error {
	// UserController_GetAllUsersFiltered: start_implement

	if c.IsAdmin(ctx.Request) {
		users, err := database.GetAllUsers()
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		var res = []*app.UserAdmin{}

		for _, user := range users {
			res = append(res, database.UserToUserAdmin(&user))
		}

		return ctx.OK(res)
	} else {
		return ctx.Forbidden(goa.ErrUnauthorized("user is not admin"))
	}

	// UserController_GetAllUsersFiltered: end_implement
}

// GetByEmail runs the get-by-email action.
func (c *UserController) GetByEmail(ctx *app.GetByEmailUserContext) error {
	// UserController_GetByEmail: start_implement

	if !c.IsAdmin(ctx.Request) {
		return ctx.NotFound(goa.ErrNotFound(ctx.RequestURI))
	}

	u, err := database.QueryUserEmail(ctx, strings.ToLower(ctx.Email))
	if err == database.ErrUserNotFound {
		return ctx.NotFound(goa.ErrNotFound(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	res := database.UserToUserAdmin(u)

	return ctx.OKAdmin(res)

	// UserController_GetByEmail: end_implement
}

// GetMany runs the get-many action.
func (c *UserController) GetMany(ctx *app.GetManyUserContext) error {
	// UserController_GetMany: start_implement

	userIDs := make([]string, 0, len(ctx.ID))
	for _, v := range ctx.ID {
		uID := v
		if uID == "" {
			return ctx.BadRequest(goa.ErrBadRequest("ID is not valid", "id", v))
		}
		userIDs = append(userIDs, uID)
	}

	users, err := database.GetUserMulti(ctx, userIDs)
	if err == database.ErrUserNotFound {
		return ctx.NotFound(goa.ErrNotFound("Unable to find all users"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if c.IsAdmin(ctx.Request) {
		res := make(app.UserAdminCollection, 0, len(users))
		for _, v := range users {
			data := database.UserToUserAdmin(v)
			res = append(res, data)
		}

		return ctx.OKAdmin(res)
	}

	res := make(app.UserCollection, 0, len(users))
	for _, v := range users {
		data := database.UserToUser(v)
		res = append(res, data)
	}

	return ctx.OK(res)

	// UserController_GetMany: end_implement
}

// GetAuths runs the getAuths action.
func (c *UserController) GetAuths(ctx *app.GetAuthsUserContext) error {
	// UserController_GetAuths: start_implement

	var uID string
	var err error

	if ctx.UserID == nil || *ctx.UserID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("User ID must be a number"))
	} else if ctx.UserID != nil && c.IsAdmin(ctx.Request) {
		uID = *ctx.UserID
	} else {
		uID = c.GetUserID(ctx.Request)
		if uID == "" {
			return ctx.Unauthorized(goa.ErrUnauthorized("User ID was not recognised."))
		}
	}

	AST := &app.AuthStatus{}

	_, err = database.QueryPasswordLoginFromID(ctx, uID)
	if err != nil {
		if err != database.ErrPasswordLoginNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else {
		AST.Standard = true
	}

	_, err = database.QueryGoogleAccountUser(ctx, uID)
	if err != nil {
		if err != database.ErrGoogleAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else {
		AST.Google = true
	}

	_, err = database.QueryFacebookAccountUser(ctx, uID)
	if err != nil {
		if err != database.ErrFacebookAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else {
		AST.Facebook = true
	}

	_, err = database.QueryLinkedinAccountUser(ctx, uID)
	if err != nil {
		if err != database.ErrLinkedinAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else {
		AST.Linkedin = true
	}

	_, err = database.QueryMicrosoftAccountUser(ctx, uID)
	if err != nil {
		if err != database.ErrMicrosoftAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else {
		AST.Microsoft = true
	}

	_, err = database.QueryTwitterAccountUser(ctx, uID)
	if err != nil {
		if err != database.ErrTwitterAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else {
		AST.Twitter = true
	}

	return ctx.OK(AST)

	// UserController_GetAuths: end_implement
}

// ResendVerifyEmail runs the resend-verify-email action.
func (c *UserController) ResendVerifyEmail(ctx *app.ResendVerifyEmailUserContext) error {
	// UserController_ResendVerifyEmail: start_implement

	uID := c.GetUserID(ctx.Request)

	u, err := database.GetUser(ctx, uID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if u.VerifiedEmail && (u.ChangingEmail == "" || u.ChangingEmail == u.Email) {
		return ctx.NotFound(goa.ErrNotFound("No email needs verifying currently"))
	}

	oldID, err := database.QueryEmailVerificationByUserID(ctx, uID)
	if err == nil {
		err = database.DeleteEmailVerification(ctx, oldID)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	} else if err != database.ErrEmailVerificationNotFound {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	eID, err := generateEmailID(u.Email)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ev := &database.EmailVerification{
		ID:          eID,
		UserID:      uID,
		TimeExpires: time.Now().Add(emailValidateExpiration),
		Email:       u.Email,
	}
	if u.VerifiedEmail && u.ChangingEmail != u.Email && u.ChangingEmail != "" {
		ev.Email = u.ChangingEmail
	}
	err = database.CreateEmailVerification(ctx, ev)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	// @TODO: Resend verification email
	resend := struct {
		UserAbout    string
		Type         string
		IsWelcome    bool
		Verification string
	}{
		UserAbout:    uID,
		Type:         "email-changed",
		IsWelcome:    !u.VerifiedEmail,
		Verification: "/verifyemail/" + eID,
	}

	subject := "Verify your email"
	toName := u.FirstName + " " + u.LastName
	toMail := u.Email
	textContent := "Go to " + secrets.URL + resend.Verification + " to verify your email"
	htmlContent := "Click <a href=\"" + secrets.URL + resend.Verification + "\">here</a> to verify your email"

	if err := email.SendMail(subject, toName, toMail, textContent, htmlContent); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))

	// UserController_ResendVerifyEmail: end_implement
}

// Retrieve runs the retrieve action.
func (c *UserController) Retrieve(ctx *app.RetrieveUserContext) error {
	// UserController_Retrieve: start_implement

	var userID string
	var err error
	if ctx.UserID != nil && *ctx.UserID != "" {
		userID = *ctx.UserID
	}

	uID := c.GetUserID(ctx.Request)

	if userID == "" {
		userID = uID
		if uID == "" {
			return ctx.Unauthorized(goa.ErrUnauthorized("Must be logged in to view own profile"))
		}
	}

	u, err := database.GetUser(ctx, userID)
	if err == database.ErrUserNotFound {
		return ctx.NotFound(goa.ErrNotFound(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if c.IsAdmin(ctx.Request) {
		res := database.UserToUserAdmin(u)
		return ctx.OKAdmin(res)
	}
	if userID == uID {
		res := database.UserToUserOwner(u)
		return ctx.OKOwner(res)
	}
	res := database.UserToUser(u)
	return ctx.OK(res)

	// UserController_Retrieve: end_implement
}

// Update runs the update action.
func (c *UserController) Update(ctx *app.UpdateUserContext) error {
	// UserController_Update: start_implement

	uID := c.GetUserID(ctx.Request)
	u, err := database.GetUser(ctx, uID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	emailChanged := false
	oldEmail := u.Email

	if ctx.Payload.Email != nil {
		if emailcheck.IsDisposableEmail(*ctx.Payload.Email) {
			return ctx.Forbidden(errDisposable("Disposable email not allowed"))
		}
		_, err = database.QueryUserEmail(ctx, *ctx.Payload.Email)
		if err == nil {
			return ctx.Forbidden(errAlreadyExists("Email is already in use"))
		} else if err != database.ErrUserNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		emailChanged = true
		if u.VerifiedEmail {
			u.ChangingEmail = *ctx.Payload.Email
		} else {
			u.Email = *ctx.Payload.Email
		}
		ctx.Payload.Email = nil
	}

	if ctx.Payload.GetNewsletter != nil && *ctx.Payload.GetNewsletter != u.GetNewsletter {
		if *ctx.Payload.GetNewsletter {
			_ = database.AddSubscriber(models.NewsletterSubscriber{
				Email:        u.Email,
				SubscribedAt: time.Now(),
				IsActive:     true,
			})
		} else {
			_ = database.RemoveSubscriber(u.Email)
		}
	}

	u = database.UserFromUserParamsMerge(ctx.Payload, u)
	if emailChanged {
		ev := &database.EmailVerification{
			UserID:      uID,
			TimeExpires: time.Now().Add(7 * 24 * time.Hour),
		}
		if u.VerifiedEmail {
			ev.Email = u.ChangingEmail
		} else {
			ev.Email = u.Email
		}
		eID, err := generateEmailID(ev.Email)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		ev.ID = eID
		err = database.CreateEmailVerification(ctx, ev)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		// @TODO: Send notification via email
		verify := struct {
			UserAbout    string
			Type         string
			IsWelcome    bool
			Verification string
		}{
			UserAbout:    uID,
			Type:         "user-updated",
			IsWelcome:    !u.VerifiedEmail,
			Verification: "/verifyemail/" + eID,
		}

		subject := "Verify your email"
		toName := u.FirstName + " " + u.LastName
		toMail := u.Email
		textContent := "Go to " + secrets.URL + verify.Verification + " to verify your email"
		htmlContent := "Click <a href=\"" + secrets.URL + verify.Verification + "\">here</a> to verify your email"

		if err := email.SendMail(subject, toName, toMail, textContent, htmlContent); err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		if !u.VerifiedEmail {
			pl, err := database.GetPasswordLogin(ctx, oldEmail)
			if err == nil {
				pl.Email = u.Email
				err = database.UpdatePasswordLogin(ctx, pl)
				if err != nil {
					return ctx.InternalServerError(goa.ErrInternal(err))
				}
				err = database.DeletePasswordLogin(ctx, oldEmail)
				if err != nil {
					return ctx.InternalServerError(goa.ErrInternal(err))
				}
			} else if err != database.ErrPasswordLoginNotFound {
				return ctx.InternalServerError(goa.ErrInternal(err))
			}
		}
	}

	err = database.UpdateUser(ctx, u)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))

	// UserController_Update: end_implement
}

func (c *UserController) UpdateAdmin(ctx *app.UpdateAdminUserContext) error {
	if c.IsAdmin(ctx.Request) {
		uID := *ctx.UID
		u, err := database.GetUser(ctx, uID)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		u.IsAdmin = *ctx.Payload.IsAdmin
		u.VerifiedEmail = *ctx.Payload.VerifiedEmail
		u.GetNewsletter = *ctx.Payload.GetNewsletter

		err = database.UpdateUser(ctx, u)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		return ctx.OK([]byte(""))

	} else {
		return ctx.Forbidden(goa.ErrUnauthorized("user is not admin"))
	}
}

func (c *UserController) UpdatePluginPermissions(ctx *app.UpdatePluginPermissionsUserContext) error {
	return nil
}

// ValidateEmail runs the validate-email action.
func (c *UserController) ValidateEmail(ctx *app.ValidateEmailUserContext) error {
	// UserController_ValidateEmail: start_implement

	// Put your logic here
	ev, err := database.GetEmailVerification(ctx, ctx.ValidateID)
	if err == database.ErrEmailVerificationNotFound {
		return ctx.NotFound([]byte("Invalid verification code, if the code was old, resend the email so a new code will be generated"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if time.Now().After(ev.TimeExpires) {
		if err := database.DeleteEmailVerification(ctx, ctx.ValidateID); err != nil {
			return err
		}
		return ctx.NotFound([]byte("Invalid verification code, if the code was old, resend the email so a new code will be generated"))
	}

	u, err := database.GetUser(ctx, ev.UserID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if !u.VerifiedEmail {
		if u.Email != ev.Email {
			return ctx.NotFound([]byte("Email is not the same as the one currently attached to this account"))
		}
		u.VerifiedEmail = true
		u.ChangingEmail = u.Email
	} else {
		if u.ChangingEmail != ev.Email {
			return ctx.NotFound([]byte("Email is not the same as the one currently attached to this account"))
		}

		pl, err := database.GetPasswordLogin(ctx, u.Email)
		if err == nil {
			pl.Email = u.ChangingEmail
			err = database.UpdatePasswordLogin(ctx, pl)
			if err != nil {
				return ctx.InternalServerError(goa.ErrInternal(err))
			}
			err = database.DeletePasswordLogin(ctx, u.Email)
			if err != nil {
				return ctx.InternalServerError(goa.ErrInternal(err))
			}
		} else if err != database.ErrPasswordLoginNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		u.Email = u.ChangingEmail
	}

	err = database.UpdateUser(ctx, u)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if err := database.DeleteEmailVerification(ctx, ctx.ValidateID); err != nil {
		return err
	}

	http.SetCookie(ctx.ResponseWriter, &http.Cookie{
		Name:   "verifyemail",
		Value:  "true",
		Path:   "/",
		MaxAge: 5 * 60, // 5 Minutes
	})
	ctx.ResponseData.Header().Set("Location", "/")
	return ctx.SeeOther()
	// UserController_ValidateEmail: end_implement
}

func createUser(ctx context.Context, u *models.User, recaptchaResponse *string, ipAddr string) (string, error) {
	if recaptchaResponse != nil {
		err := validateRecaptcha(ctx, *recaptchaResponse, ipAddr)
		if err != nil {
			return "", err
		}
	}

	uID, err := database.CreateUser(ctx, u)
	if err != nil {
		return "", err
	}

	if u.GetNewsletter {
		_ = database.AddSubscriber(models.NewsletterSubscriber{
			Email:        u.Email,
			SubscribedAt: time.Now(),
			IsActive:     true,
		})
	}

	eID, err := generateEmailID(u.Email)
	if err != nil {
		return "", err
	}
	ev := &database.EmailVerification{
		ID:          eID,
		UserID:      uID,
		TimeExpires: time.Now().Add(emailValidateExpiration),
		Email:       u.Email,
	}
	err = database.CreateEmailVerification(ctx, ev)
	if err != nil {
		return "", err
	}

	// @TODO: Use this struct to send actual email for verification
	verify := struct {
		UserAbout    string
		Type         string
		IsWelcome    bool
		Verification string
	}{
		UserAbout:    uID,
		Type:         "user-created",
		IsWelcome:    true,
		Verification: "/verifyemail/" + eID,
	}

	subject := "Verify your email"
	toName := u.FirstName + " " + u.LastName
	toMail := u.Email
	textContent := "Go to " + secrets.URL + verify.Verification + " to verify your email"
	htmlContent := "Click <a href=\"" + secrets.URL + verify.Verification + "\">here</a> to verify your email"

	if err := email.SendMail(subject, toName, toMail, textContent, htmlContent); err != nil {
		return "", err
	}

	return uID, nil
}

func validateRecaptcha(ctx context.Context, recaptchaResponse, ipAddr string) error {
	//v := url.Values{}
	//v.Set("secret", secrets.RecaptchaSecret)
	//v.Set("response", recaptchaResponse)
	//v.Set("remoteip", ipAddr)
	//c := urlfetch.Client(ctx)
	//recapRes, err := c.PostForm("https://www.google.com/recaptcha/api/siteverify", v)
	//if err != nil {
	//	return err
	//}
	//defer recapRes.Body.Close()
	//
	//var recapResData struct {
	//	Success    bool
	//	ErrorCodes []string `json:"error-codes"`
	//}
	//err = json.NewDecoder(recapRes.Body).Decode(&recapResData)
	//if err != nil {
	//	return err
	//}
	//if !recapResData.Success {
	//	return ErrInvalidRecaptcha
	//}
	return nil
}

func generateEmailID(emailAddr string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, emailAddr)
	if err != nil {
		return "", err
	}
	_, err = io.WriteString(h, time.Now().String())
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func getNumLoginMethods(ctx context.Context, userID string) int64 {
	getAccountItems := []func(context.Context, string) (string, error){
		database.QueryGoogleAccountUser,
		database.QueryFacebookAccountUser,
		database.QueryPasswordLoginFromID,
	}
	var wg sync.WaitGroup
	var numLogins int64
	wg.Add(len(getAccountItems))
	for _, v := range getAccountItems {
		go func(command func(context.Context, string) (string, error)) {
			_, err := command(ctx, userID)
			if err == nil {
				atomic.AddInt64(&numLogins, 1)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	return numLogins
}

type signer string

func (s signer) Sign(req *http.Request) error {
	req.Header.Set("Authorization", string(s))
	return nil
}
