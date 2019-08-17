package auth

import (
	"fmt"
	"github.com/alioygur/is"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/auth/models"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/email"
	"gigglesearch.org/giggle-auth/utils/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/goadesign/goa"
	"github.com/gofrs/uuid"
	"github.com/simukti/emailcheck"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/url"
	"strings"
	"time"
)

var errEmailExists = goa.ErrBadRequest("email already exists")

// PasswordAuthController implements the password-auth resource.
type PasswordAuthController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

// NewPasswordAuthController creates a password-auth controller.
func NewPasswordAuthController(service *goa.Service, jwtSec auth.JWTSecurity, sesCont *SessionController) *PasswordAuthController {
	return &PasswordAuthController{
		Controller:        service.NewController("PasswordAuthController"),
		JWTSecurity:       jwtSec,
		sessionController: sesCont,
	}
}

// ChangePassword runs the change-password action.
func (c *PasswordAuthController) ChangePassword(ctx *app.ChangePasswordPasswordAuthContext) error {
	// PasswordAuthController_ChangePassword: start_implement

	uID := c.GetUserID(ctx.Request)

	u, err := database.GetUser(ctx, uID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	passl, err := database.GetPasswordLogin(ctx, strings.ToLower(u.Email))
	if err == database.ErrPasswordLoginNotFound {
		fmt.Println(database.ErrPasswordLoginNotFound)
		return ctx.BadRequest(goa.ErrBadRequest("incorrect old password"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(passl.Password), []byte(ctx.Payload.OldPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(bcrypt.ErrMismatchedHashAndPassword)
		return ctx.BadRequest(goa.ErrBadRequest("incorrect old password"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	cryptPass, err := bcrypt.GenerateFromPassword([]byte(ctx.Payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	newP := database.PasswordLogin{
		Email:    u.Email,
		UserID:   u.ID.Hex(),
		Password: string(cryptPass),
	}
	err = database.UpdatePasswordLogin(ctx, &newP)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	sesID := c.GetSessionFromAuth(ctx.Request)
	err = logoutAllSessionsBut(ctx, uID, sesID)
	if err != nil {
		log.Warning(ctx, "Unable to logout of other sessions when changing password, uID=%s, err=%v", uID, err)
	}

	return ctx.OK([]byte(""))

	// PasswordAuthController_ChangePassword: end_implement
}

// ConfirmReset runs the confirm-reset action.
func (c *PasswordAuthController) ConfirmReset(ctx *app.ConfirmResetPasswordAuthContext) error {
	// PasswordAuthController_ConfirmReset: start_implement

	rp, err := database.GetResetPassword(ctx, ctx.Payload.UserID)
	if err == database.ErrResetPasswordNotFound {
		return ctx.Forbidden(errBadReset("Invalid reset code"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	fmt.Println(rp)

	if time.Now().After(rp.TimeExpires) {
		database.DeleteResetPassword(ctx, rp.UserID)
		return ctx.Forbidden(errBadReset("Invalid reset code"))
	}

	fmt.Println(rp.ID.String(), ctx.Payload.ResetCode)

	if rp.ID.String() != ctx.Payload.ResetCode {
		return ctx.Forbidden(errBadReset("Invalid reset code"))
	}

	u, err := database.GetUser(ctx, rp.UserID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	cryptPass, err := bcrypt.GenerateFromPassword([]byte(ctx.Payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = database.UpdatePasswordLogin(ctx, &database.PasswordLogin{
		Email:    u.Email,
		UserID:   u.ID.Hex(),
		Password: string(cryptPass),
	})
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = logoutAllSessionsBut(ctx, rp.UserID, "")
	if err != nil {
		log.Warning(ctx, "Unable to logout of all sessions when resetting password: %v", err)
	}

	return ctx.OK([]byte(""))

	// PasswordAuthController_ConfirmReset: end_implement
}

// Login runs the login action.
func (c *PasswordAuthController) Login(ctx *app.LoginPasswordAuthContext) error {
	// PasswordAuthController_Login: start_implement

	passl, err := database.GetPasswordLogin(ctx, strings.ToLower(ctx.Payload.Email))
	if err == database.ErrPasswordLoginNotFound {
		fmt.Println(database.ErrPasswordLoginNotFound)
		return ctx.Unauthorized(goa.ErrUnauthorized("Email or password does not match"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(passl.Password), []byte(ctx.Payload.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(bcrypt.ErrMismatchedHashAndPassword)
		return ctx.Unauthorized(goa.ErrUnauthorized("Email or password does not match"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	u, err := database.GetUser(ctx, passl.UserID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *u, ctx.Token)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
	return ctx.OK(database.UserToUser(u))

	// PasswordAuthController_Login: end_implement
}

// Register runs the register action.
func (c *PasswordAuthController) Register(ctx *app.RegisterPasswordAuthContext) error {
	// PasswordAuthController_Register: start_implement

	if !is.Email(strings.ToLower(ctx.Payload.Email)) {
		return ctx.BadRequest(errEmailExists)
	}
	if emailcheck.IsDisposableEmail(strings.ToLower(ctx.Payload.Email)) {
		return ctx.Forbidden(goa.ErrBadRequest("disposable email not allowed"))
	}

	if emailExist := CheckEmailExists(ctx.Payload.Email); emailExist {
		return ctx.Forbidden(errEmailExists)
	}

	_, _, err := net.SplitHostPort(ctx.RequestData.RemoteAddr)
	if err != nil {
		_ = ctx.RequestData.RemoteAddr
	}

	payload := ctx.Payload

	cryptPass, err := bcrypt.GenerateFromPassword([]byte(ctx.Payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	newU := models.User{
		ID:            bson.NewObjectId(),
		Email:         payload.Email,
		FirstName:     payload.FirstName,
		LastName:      payload.LastName,
		Category:      payload.Category,
		Password:      string(cryptPass),
		VerifiedEmail: false,
		IsAdmin:       false,
		GetNewsletter: false,
	}

	if err := models.UsersCollection.Insert(&newU); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	passl := models.PasswordLogin{
		Email:    strings.ToLower(ctx.Payload.Email),
		UserID:   newU.ID.Hex(),
		Password: string(cryptPass),
	}

	fmt.Println(passl)

	if err := models.PasswordLoginCollection.Insert(&passl); err != nil {
		fmt.Printf("%v", err)
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, newU, nil)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
	return ctx.OK(database.UserToUser(&newU))

	// PasswordAuthController_Register: end_implement
}

// Remove runs the remove action.
func (c *PasswordAuthController) Remove(ctx *app.RemovePasswordAuthContext) error {
	// PasswordAuthController_Remove: start_implement

	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot remove password if it is the only way to login"))
	}

	pID, err := database.QueryPasswordLoginFromID(ctx, uID)
	if err == database.ErrPasswordLoginNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account does not have a password"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeletePasswordLogin(ctx, pID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))

	// PasswordAuthController_Remove: end_implement
}

// Reset runs the reset action.
func (c *PasswordAuthController) Reset(ctx *app.ResetPasswordAuthContext) error {
	// PasswordAuthController_Reset: start_implement

	u, err := database.QueryUserEmail(ctx, strings.ToLower(ctx.Email))
	if err == database.ErrUserNotFound {
		return ctx.OK([]byte(""))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	var id uuid.UUID
	rp, err := database.GetResetPassword(ctx, u.ID.Hex())
	if err == nil {
		id = rp.ID
	} else if err == database.ErrResetPasswordNotFound {
		id, err = uuid.NewV4()
		if err != nil {
			return err
		}
	} else {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = database.CreateResetPassword(ctx, &database.ResetPassword{
		UserID:      u.ID.Hex(),
		ID:          id,
		TimeExpires: time.Now().Add(120 * time.Minute),
	})
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	resetPass := struct {
		UserAbout    string
		Type         string
		Verification string
	}{
		UserAbout:    u.ID.Hex(),
		Type:         "password-reset",
		Verification: "?code=" + url.QueryEscape(id.String()) + "&uid=" + url.QueryEscape(u.ID.Hex()),
	}

	textContent := "Go to " + ctx.RedirectURL + resetPass.Verification + " to reset your password"
	htmlContent := "Click <a href=\"" + ctx.RedirectURL + resetPass.Verification + "\">here</a> to reset your password"
	toName := u.FirstName + " " + u.LastName
	toMail := u.Email

	if err := email.SendMail("Reset your Password", toName, toMail, textContent, htmlContent); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))

	// PasswordAuthController_Reset: end_implement
}

func CheckEmailExists(email string) bool {
	if count, err := models.UsersCollection.Find(bson.M{"email": email}).Count(); err == mgo.ErrNotFound {
		return false
	} else if err != nil {
		return true
	} else if count > 0 {
		return true
	}

	return false
}
