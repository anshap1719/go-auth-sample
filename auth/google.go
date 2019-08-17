package auth

import (
	"encoding/json"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/log"
	"gigglesearch.org/giggle-auth/utils/secrets"
	"github.com/goadesign/goa"
	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
	ght "golang.org/x/oauth2/google"
	"net"
	"time"
)

// GoogleController implements the google resource.
type GoogleController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

const (
	googleRegister = iota
	googleLogin
	googleAttach
)

const (
	googleConnectionExpire = 30 * time.Minute
	googleRegisterExpire   = time.Hour
)

var googleConf = &oauth2.Config{
	ClientID:     secrets.GoogleID,
	ClientSecret: secrets.GoogleSecret,
	Scopes: []string{
		"email",
		"profile",
	},
	Endpoint:    ght.Endpoint,
	RedirectURL: "https://localhost:4000/social",
}

// NewGoogleController creates a google controller.
func NewGoogleController(service *goa.Service, jwtSec auth.JWTSecurity, sesCont *SessionController) *GoogleController {
	return &GoogleController{
		Controller:        service.NewController("GoogleController"),
		JWTSecurity:       jwtSec,
		sessionController: sesCont,
	}
}

// AttachToAccount runs the attach-to-account action.
func (c *GoogleController) AttachToAccount(ctx *app.AttachToAccountGoogleContext) error {
	// GoogleController_AttachToAccount: start_implement

	gc := &database.GoogleConnection{
		TimeCreated: time.Now(),
		Purpose:     googleAttach,
	}
	state, err := database.CreateGoogleConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(googleConf.AuthCodeURL(state.String())))

	// GoogleController_AttachToAccount: end_implement
}

// DetachFromAccount runs the detach-from-account action.
func (c *GoogleController) DetachFromAccount(ctx *app.DetachFromAccountGoogleContext) error {
	// GoogleController_DetachFromAccount: start_implement

	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot detach last login method"))
	}

	gID, err := database.QueryGoogleAccountUser(ctx, uID)
	if err == database.ErrGoogleAccountNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account is not connected to Google"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteGoogleAccount(ctx, gID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))

	// GoogleController_DetachFromAccount: end_implement
}

// Login runs the login action.
func (c *GoogleController) Login(ctx *app.LoginGoogleContext) error {
	// GoogleController_Login: start_implement

	var mt uuid.UUID
	if ctx.Token != nil {
		mt = *ctx.Token
	}

	gc := &database.GoogleConnection{
		TimeCreated: time.Now(),
		Purpose:     googleLogin,
		MergeToken:  mt,
	}

	state, err := database.CreateGoogleConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(googleConf.AuthCodeURL(state.String())))

	// GoogleController_Login: end_implement
}

type GoogleUser struct {
	ID        *string `json:"id"`
	FirstName string  `json:"given_name"`
	LastName  string  `json:"family_name"`
	Email     *string `json:"email"`
}

func (g *GoogleUser) GetEmail() string {
	if g.Email == nil {
		return ""
	}
	return *g.Email
}

// Receive runs the receive action.
func (c *GoogleController) Receive(ctx *app.ReceiveGoogleContext) error {
	// GoogleController_Receive: start_implement

	gc, err := database.GetGoogleConnection(*ctx, ctx.State)
	if err == database.ErrGoogleConnectionNotFound {
		return ctx.BadRequest(goa.ErrBadRequest("Google connection must be created with other API methods"))
	} else if err != nil {

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = database.DeleteGoogleConnection(ctx, ctx.State)
	if err != nil {
		log.Warning(ctx, "Unable to delete Google connection, state=%s", ctx.State)
	}

	if gc.TimeCreated.Add(googleConnectionExpire).Before(time.Now()) {
		return ctx.BadRequest(goa.ErrBadRequest("Google connection must be created with other API methods"))
	}

	token, err := googleConf.Exchange(ctx, ctx.Code)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	tokenClient := googleConf.Client(ctx, token)

	resp, err := tokenClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid Google data"))
	}

	googleUser := &GoogleUser{}

	err = json.NewDecoder(resp.Body).Decode(googleUser)
	if err != nil {

		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	gID := googleUser.GetEmail()
	if gID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Unable to get Google user ID"))
	}

	switch gc.Purpose {
	case googleRegister:
		_, err := database.GetGoogleAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Google account is already attached to an account"))
		} else if err != database.ErrGoogleAccountNotFound {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		gr := &database.GoogleRegister{
			GoogleEmail: gID,
			TimeCreated: time.Now(),
		}
		regID, err := database.CreateGoogleRegister(ctx, gr)
		if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		grm := &app.GoogleRegisterMedia{
			OauthKey:  regID,
			FirstName: googleUser.FirstName,
			LastName:  googleUser.LastName,
			Email:     googleUser.GetEmail(),
		}
		return ctx.OK(grm)
	case googleLogin:
		account, err := database.GetGoogleAccount(ctx, gID)
		if err == database.ErrGoogleAccountNotFound {
			return ctx.BadRequest(goa.ErrBadRequest("No account associated with that Google account"))
		} else if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		u, err := database.GetUser(ctx, account.UserID)
		if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		var mt *uuid.UUID

		if gc.MergeToken.String() == "" {
			mt = &gc.MergeToken
		}

		sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *u, mt)
		if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		ctx.ResponseData.Header().Set("X-Session", sesToken)
		ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)

		return ctx.OK(database.UserToUser(u))
	case googleAttach:
		_, err := database.GetGoogleAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Google account is already attached to an account"))
		} else if err != database.ErrGoogleAccountNotFound {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		uID := c.GetUserID(ctx.Request)
		if uID == "" {
			return ctx.Unauthorized(goa.ErrUnauthorized("You must be logged in"))
		}
		_, err = database.GetUser(ctx, uID)
		if err == database.ErrUserNotFound {
			log.Critical(ctx, "Unable to get user account, userID=%s", uID)
			return ctx.InternalServerError(goa.ErrInternal("Unable to find user account"))
		} else if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		account := &database.GoogleAccount{
			GoogleEmail: gID,
			UserID:      uID,
		}
		err = database.CreateGoogleAccount(ctx, account)
		if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	default:
		log.Critical(ctx, "Bad Google receive type")
		return ctx.InternalServerError(goa.ErrInternal("Invalid Google connection type"))
	}

	return ctx.OK(nil)

	// GoogleController_Receive: end_implement
}

// Register runs the register action.
func (c *GoogleController) Register(ctx *app.RegisterGoogleContext) error {
	// GoogleController_Register: start_implement

	gr, err := database.GetGoogleRegister(ctx, ctx.Payload.OauthKey)
	if err == database.ErrGoogleRegisterNotFound {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if gr.TimeCreated.Add(googleRegisterExpire).Before(time.Now()) {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	}
	_, err = database.GetGoogleAccount(ctx, gr.GoogleEmail)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This Google account is already attached to an account"))
	} else if err != database.ErrGoogleAccountNotFound {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	_, err = database.QueryUserEmail(ctx, ctx.Payload.Email)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This email is already in use", "email", ctx.Payload.Email))
	} else if err != database.ErrUserNotFound {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ipAddr, _, err := net.SplitHostPort(ctx.RequestData.RemoteAddr)
	if err != nil {
		ipAddr = ctx.RequestData.RemoteAddr
	}

	newU := database.UserFromGoogleRegisterParams(ctx.Payload)
	uID, err := createUser(ctx, newU, &ctx.Payload.GRecaptchaResponse, ipAddr)
	if err == ErrInvalidRecaptcha {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	account := &database.GoogleAccount{
		GoogleEmail: gr.GoogleEmail,
		UserID:      uID,
	}
	err = database.CreateGoogleAccount(ctx, account)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteGoogleRegister(ctx, ctx.Payload.OauthKey)
	if err != nil {
		log.Warning(ctx, "Unable to delete Google registration progress, Google Key=%v", ctx.Payload.OauthKey)
	}

	sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *newU, nil)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)

	newU.Email = gr.GoogleEmail

	return ctx.OK(database.UserToUser(newU))

	// GoogleController_Register: end_implement
}

// RegisterURL runs the register-url action.
func (c *GoogleController) RegisterURL(ctx *app.RegisterURLGoogleContext) error {
	// GoogleController_RegisterURL: start_implement

	gc := &database.GoogleConnection{
		TimeCreated: time.Now(),
		Purpose:     googleRegister,
	}
	state, err := database.CreateGoogleConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(googleConf.AuthCodeURL(state.String())))

	// GoogleController_RegisterURL: end_implement
}
