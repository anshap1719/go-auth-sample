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
	ght "golang.org/x/oauth2/linkedin"
	"net"
	"time"
)

// LinkedinController implements the linkedin resource.
type LinkedinController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

const (
	linkedinRegister = iota
	linkedinLogin
	linkedinAttach
)

const (
	linkedinConnectionExpire = 30 * time.Minute
	linkedinRegisterExpire   = time.Hour
)

var linkedinConf = &oauth2.Config{
	ClientID:     secrets.LinkedInID,
	ClientSecret: secrets.LinkedInSecret,
	Scopes: []string{
		"r_emailaddress",
		"r_basicprofile",
	},
	Endpoint:    ght.Endpoint,
	RedirectURL: "https://localhost:4000/ln_social",
}

// NewLinkedinController creates a linkedin controller.
func NewLinkedinController(service *goa.Service, jwtSec auth.JWTSecurity, sesCont *SessionController) *LinkedinController {
	return &LinkedinController{
		Controller:        service.NewController("LinkedinController"),
		JWTSecurity:       jwtSec,
		sessionController: sesCont,
	}
}

// AttachToAccount runs the attach-to-account action.
func (c *LinkedinController) AttachToAccount(ctx *app.AttachToAccountLinkedinContext) error {
	gc := &database.LinkedinConnection{
		TimeCreated: time.Now(),
		Purpose:     linkedinAttach,
	}
	state, err := database.CreateLinkedinConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(linkedinConf.AuthCodeURL(state.String())))
}

// DetachFromAccount runs the detach-from-account action.
func (c *LinkedinController) DetachFromAccount(ctx *app.DetachFromAccountLinkedinContext) error {
	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot detach last login method"))
	}

	gID, err := database.QueryLinkedinAccountUser(ctx, uID)
	if err == database.ErrLinkedinAccountNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account is not connected to Linkedin"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteLinkedinAccount(ctx, gID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))
}

// Login runs the login action.
func (c *LinkedinController) Login(ctx *app.LoginLinkedinContext) error {
	var mt uuid.UUID
	if ctx.Token != nil {
		mt = *ctx.Token
	}
	gc := &database.LinkedinConnection{
		TimeCreated: time.Now(),
		Purpose:     linkedinLogin,
		MergeToken:  mt,
	}
	state, err := database.CreateLinkedinConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(linkedinConf.AuthCodeURL(state.String())))
}

type LinkedinUser struct {
	ID        string `json:"id" bson:"id"`
	FirstName string `json:"firstName" bson:"first_name"`
	LastName  string `json:"lastName" bson:"last_name"`
	Email     string `json:"emailAddress" bson:"email_address"`
}

// Receive runs the receive action.
func (c *LinkedinController) Receive(ctx *app.ReceiveLinkedinContext) error {
	gc, err := database.GetLinkedinConnection(ctx, ctx.State)
	if err == database.ErrLinkedinConnectionNotFound {
		return ctx.BadRequest(goa.ErrBadRequest("Linkedin connection must be created with other API methods"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteLinkedinConnection(ctx, ctx.State)
	if err != nil {
		log.Warning(ctx, "Unable to delete Linkedin connection, state=%s", ctx.State)
	}
	if gc.TimeCreated.Add(linkedinConnectionExpire).Before(time.Now()) {
		return ctx.BadRequest(goa.ErrBadRequest("Linkedin connection must be created with other API methods"))
	}

	token, err := linkedinConf.Exchange(ctx, ctx.Code)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	tokenClient := linkedinConf.Client(ctx, token)
	resp, err := tokenClient.Get("https://api.linkedin.com/v1/people/~:(id,first-name,email-address,last-name)?format=json")
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid Linkedin data"))
	}
	linkedinUser := &LinkedinUser{}
	err = json.NewDecoder(resp.Body).Decode(linkedinUser)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	gID := linkedinUser.Email
	if gID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Unable to get Linkedin user ID"))
	}

	switch gc.Purpose {
	case linkedinRegister:
		_, err := database.GetLinkedinAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Linkedin account is already attached to an account"))
		} else if err != database.ErrLinkedinAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		gr := &database.LinkedinRegister{
			LinkedinEmail: gID,
			TimeCreated:   time.Now(),
		}
		regID, err := database.CreateLinkedinRegister(ctx, gr)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		grm := &app.LinkedinRegisterMedia{
			OauthKey:  regID,
			FirstName: linkedinUser.FirstName,
			LastName:  linkedinUser.LastName,
			Email:     linkedinUser.Email,
		}
		return ctx.OK(grm)
	case linkedinLogin:
		account, err := database.GetLinkedinAccount(ctx, gID)
		if err == database.ErrLinkedinAccountNotFound {
			return ctx.BadRequest(goa.ErrBadRequest("No account associated with that Linkedin account"))
		} else if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		u, err := database.GetUser(ctx, account.UserID)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		var mt *uuid.UUID

		if gc.MergeToken.String() != "" {
			mt = &gc.MergeToken
		}

		sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *u, mt)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		ctx.ResponseData.Header().Set("X-Session", sesToken)
		ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
		return ctx.OK(database.UserToUser(u))
	case linkedinAttach:
		_, err := database.GetLinkedinAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Linkedin account is already attached to an account"))
		} else if err != database.ErrLinkedinAccountNotFound {
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

		account := &database.LinkedinAccount{
			LinkedinEmail: gID,
			UserID:        uID,
		}
		err = database.CreateLinkedinAccount(ctx, account)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	default:
		log.Critical(ctx, "Bad Linkedin receive type")
		return ctx.InternalServerError(goa.ErrInternal("Invalid Linkedin connection type"))
	}

	return ctx.OK(nil)
}

// Register runs the register action.
func (c *LinkedinController) Register(ctx *app.RegisterLinkedinContext) error {
	gr, err := database.GetLinkedinRegister(ctx, ctx.Payload.OauthKey)
	if err == database.ErrLinkedinRegisterNotFound {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if gr.TimeCreated.Add(linkedinRegisterExpire).Before(time.Now()) {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	}
	_, err = database.GetLinkedinAccount(ctx, gr.LinkedinEmail)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This Linkedin account is already attached to an account"))
	} else if err != database.ErrLinkedinAccountNotFound {
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

	newU := database.UserFromLinkedinRegisterParams(ctx.Payload)
	uID, err := createUser(ctx, newU, &ctx.Payload.GRecaptchaResponse, ipAddr)
	if err == ErrInvalidRecaptcha {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	account := &database.LinkedinAccount{
		LinkedinEmail: gr.LinkedinEmail,
		UserID:        uID,
	}
	err = database.CreateLinkedinAccount(ctx, account)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteLinkedinRegister(ctx, ctx.Payload.OauthKey)
	if err != nil {
		log.Warning(ctx, "Unable to delete Linkedin registration progress, Linkedin Key=%v", ctx.Payload.OauthKey)
	}

	sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *newU, nil)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
	return ctx.OK(database.UserToUser(newU))
}

// RegisterURL runs the register-url action.
func (c *LinkedinController) RegisterURL(ctx *app.RegisterURLLinkedinContext) error {
	gc := &database.LinkedinConnection{
		TimeCreated: time.Now(),
		Purpose:     linkedinRegister,
	}
	state, err := database.CreateLinkedinConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(linkedinConf.AuthCodeURL(state.String())))
}
