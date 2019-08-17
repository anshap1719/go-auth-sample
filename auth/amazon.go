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
	ght "golang.org/x/oauth2/amazon"
	"net"
	"strings"
	"time"
)

// AmazonController implements the amazon resource.
type AmazonController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

const (
	amazonRegister = iota
	amazonLogin
	amazonAttach
)

const (
	amazonConnectionExpire = 30 * time.Minute
	amazonRegisterExpire   = time.Hour
)

var amazonConf = &oauth2.Config{
	ClientID:     secrets.AmazonID,
	ClientSecret: secrets.AmazonSecret,
	Scopes: []string{
		"profile",
	},
	Endpoint:    ght.Endpoint,
	RedirectURL: "https://localhost:4000/am_social",
}

// NewAmazonController creates a amazon controller.
func NewAmazonController(service *goa.Service, jwtSec auth.JWTSecurity, sesCont *SessionController) *AmazonController {
	return &AmazonController{
		Controller:        service.NewController("AmazonController"),
		JWTSecurity:       jwtSec,
		sessionController: sesCont,
	}
}

// AttachToAccount runs the attach-to-account action.
func (c *AmazonController) AttachToAccount(ctx *app.AttachToAccountAmazonContext) error {
	// AmazonController_AttachToAccount: start_implement

	gc := &database.AmazonConnection{
		TimeCreated: time.Now(),
		Purpose:     amazonAttach,
	}
	state, err := database.CreateAmazonConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(amazonConf.AuthCodeURL(state.String())))

	// AmazonController_AttachToAccount: end_implement
}

// DetachFromAccount runs the detach-from-account action.
func (c *AmazonController) DetachFromAccount(ctx *app.DetachFromAccountAmazonContext) error {
	// AmazonController_DetachFromAccount: start_implement

	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot detach last login method"))
	}

	gID, err := database.QueryAmazonAccountUser(ctx, uID)
	if err == database.ErrAmazonAccountNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account is not connected to Amazon"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteAmazonAccount(ctx, gID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))

	// AmazonController_DetachFromAccount: end_implement
}

// Login runs the login action.
func (c *AmazonController) Login(ctx *app.LoginAmazonContext) error {
	// AmazonController_Login: start_implement

	var mt uuid.UUID
	if ctx.Token != nil {
		mt = *ctx.Token
	}

	gc := &database.AmazonConnection{
		TimeCreated: time.Now(),
		Purpose:     amazonLogin,
		MergeToken:  mt,
	}

	state, err := database.CreateAmazonConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(amazonConf.AuthCodeURL(state.String())))

	// AmazonController_Login: end_implement
}

type AmazonUser struct {
	ID    *string `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (g *AmazonUser) GetEmail() string {
	if g.Email == nil {
		return ""
	}
	return *g.Email
}

// Receive runs the receive action.
func (c *AmazonController) Receive(ctx *app.ReceiveAmazonContext) error {
	// AmazonController_Receive: start_implement

	gc, err := database.GetAmazonConnection(*ctx, ctx.State)
	if err == database.ErrAmazonConnectionNotFound {
		return ctx.BadRequest(goa.ErrBadRequest("Amazon connection must be created with other API methods"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = database.DeleteAmazonConnection(ctx, ctx.State)
	if err != nil {
		log.Warning(ctx, "Unable to delete Amazon connection, state=%s", ctx.State)
	}

	if gc.TimeCreated.Add(amazonConnectionExpire).Before(time.Now()) {
		return ctx.BadRequest(goa.ErrBadRequest("Amazon connection must be created with other API methods"))
	}

	token, err := amazonConf.Exchange(ctx, ctx.Code)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	tokenClient := amazonConf.Client(ctx, token)

	resp, err := tokenClient.Get("https://api.amazon.com/user/profile")
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid Amazon data"))
	}

	amazonUser := &AmazonUser{}

	err = json.NewDecoder(resp.Body).Decode(amazonUser)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	gID := amazonUser.GetEmail()
	if gID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Unable to get Amazon user ID"))
	}

	switch gc.Purpose {
	case amazonRegister:
		_, err := database.GetAmazonAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Amazon account is already attached to an account"))
		} else if err != database.ErrAmazonAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		gr := &database.AmazonRegister{
			AmazonEmail: gID,
			TimeCreated: time.Now(),
		}

		regID, err := database.CreateAmazonRegister(ctx, gr)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		grm := &app.AmazonRegisterMedia{
			OauthKey:  regID,
			FirstName: strings.Split(*amazonUser.Name, " ")[0],
			LastName:  strings.Split(*amazonUser.Name, " ")[1],
			Email:     amazonUser.GetEmail(),
		}
		return ctx.OK(grm)
	case amazonLogin:
		account, err := database.GetAmazonAccount(ctx, gID)

		if err == database.ErrAmazonAccountNotFound {
			return ctx.BadRequest(goa.ErrBadRequest("No account associated with that Amazon account"))
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
	case amazonAttach:
		_, err := database.GetAmazonAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Amazon account is already attached to an account"))
		} else if err != database.ErrAmazonAccountNotFound {
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

		account := &database.AmazonAccount{
			AmazonEmail: gID,
			UserID:      uID,
		}
		err = database.CreateAmazonAccount(ctx, account)
		if err != nil {

			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	default:
		log.Critical(ctx, "Bad Amazon receive type")
		return ctx.InternalServerError(goa.ErrInternal("Invalid Amazon connection type"))
	}

	return ctx.OK(nil)

	// AmazonController_Receive: end_implement
}

// Register runs the register action.
func (c *AmazonController) Register(ctx *app.RegisterAmazonContext) error {
	// AmazonController_Register: start_implement

	gr, err := database.GetAmazonRegister(ctx, ctx.Payload.OauthKey)
	if err == database.ErrAmazonRegisterNotFound {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if gr.TimeCreated.Add(amazonRegisterExpire).Before(time.Now()) {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	}
	_, err = database.GetAmazonAccount(ctx, gr.AmazonEmail)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This Amazon account is already attached to an account"))
	} else if err != database.ErrAmazonAccountNotFound {
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

	newU := database.UserFromAmazonRegisterParams(ctx.Payload)
	uID, err := createUser(ctx, newU, &ctx.Payload.GRecaptchaResponse, ipAddr)
	if err == ErrInvalidRecaptcha {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	account := &database.AmazonAccount{
		AmazonEmail: gr.AmazonEmail,
		UserID:      uID,
	}
	err = database.CreateAmazonAccount(ctx, account)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteAmazonRegister(ctx, ctx.Payload.OauthKey)
	if err != nil {
		log.Warning(ctx, "Unable to delete Amazon registration progress, Amazon Key=%v", ctx.Payload.OauthKey)
	}

	sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *newU, nil)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
	return ctx.OK(database.UserToUser(newU))

	// AmazonController_Register: end_implement
}

// RegisterURL runs the register-url action.
func (c *AmazonController) RegisterURL(ctx *app.RegisterURLAmazonContext) error {
	// AmazonController_RegisterURL: start_implement

	gc := &database.AmazonConnection{
		TimeCreated: time.Now(),
		Purpose:     amazonRegister,
	}
	state, err := database.CreateAmazonConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(amazonConf.AuthCodeURL(state.String())))

	// AmazonController_RegisterURL: end_implement
}
