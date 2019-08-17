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
	ght "golang.org/x/oauth2/microsoft"
	"net"
	"time"
)

// MicrosoftController implements the microsoft resource.
type MicrosoftController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

const (
	microsoftRegister = iota
	microsoftLogin
	microsoftAttach
)

const (
	microsoftConnectionExpire = 30 * time.Minute
	microsoftRegisterExpire   = time.Hour
)

var microsoftConf = &oauth2.Config{
	ClientID:     secrets.MicrosoftID,
	ClientSecret: secrets.MicrosoftSecret,
	Scopes: []string{
		"wl.signin",
		"wl.emails",
	},
	Endpoint:    ght.LiveConnectEndpoint,
	RedirectURL: "https://localhost:4000/ms_social",
}

// NewMicrosoftController creates a microsoft controller.
func NewMicrosoftController(service *goa.Service, jwtSec auth.JWTSecurity, sesCont *SessionController) *MicrosoftController {
	return &MicrosoftController{
		Controller:        service.NewController("MicrosoftController"),
		JWTSecurity:       jwtSec,
		sessionController: sesCont,
	}
}

// AttachToAccount runs the attach-to-account action.
func (c *MicrosoftController) AttachToAccount(ctx *app.AttachToAccountMicrosoftContext) error {
	gc := &database.MicrosoftConnection{
		TimeCreated: time.Now(),
		Purpose:     microsoftAttach,
	}
	state, err := database.CreateMicrosoftConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(microsoftConf.AuthCodeURL(state.String())))
}

// DetachFromAccount runs the detach-from-account action.
func (c *MicrosoftController) DetachFromAccount(ctx *app.DetachFromAccountMicrosoftContext) error {
	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot detach last login method"))
	}

	gID, err := database.QueryMicrosoftAccountUser(ctx, uID)
	if err == database.ErrMicrosoftAccountNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account is not connected to Microsoft"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteMicrosoftAccount(ctx, gID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))
}

// Login runs the login action.
func (c *MicrosoftController) Login(ctx *app.LoginMicrosoftContext) error {
	var mt uuid.UUID
	if ctx.Token != nil {
		mt = *ctx.Token
	}
	gc := &database.MicrosoftConnection{
		TimeCreated: time.Now(),
		Purpose:     microsoftLogin,
		MergeToken:  mt,
	}
	state, err := database.CreateMicrosoftConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(microsoftConf.AuthCodeURL(state.String())))
}

type MicrosoftUser struct {
	ID        string `json:"id" bson:"id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     struct {
		Address string `json:"account" bson:"account"`
	} `json:"emails" bson:"emails"`
}

// Receive runs the receive action.
func (c *MicrosoftController) Receive(ctx *app.ReceiveMicrosoftContext) error {
	gc, err := database.GetMicrosoftConnection(ctx, ctx.State)
	if err == database.ErrMicrosoftConnectionNotFound {
		return ctx.BadRequest(goa.ErrBadRequest("Microsoft connection must be created with other API methods"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteMicrosoftConnection(ctx, ctx.State)
	if err != nil {
		log.Warning(ctx, "Unable to delete Microsoft connection, state=%s", ctx.State)
	}
	if gc.TimeCreated.Add(microsoftConnectionExpire).Before(time.Now()) {
		return ctx.BadRequest(goa.ErrBadRequest("Microsoft connection must be created with other API methods"))
	}

	token, err := microsoftConf.Exchange(ctx, ctx.Code)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	tokenClient := microsoftConf.Client(ctx, token)
	resp, err := tokenClient.Get("https://apis.live.net/v5.0/me")
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid Microsoft data"))
	}
	microsoftUser := &MicrosoftUser{}
	err = json.NewDecoder(resp.Body).Decode(microsoftUser)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	// cli := microsoft.NewClient(tokenClient)
	// microsoftUser, _, err := cli.Users.Get(ctx, "")
	// if err != nil {
	// 	return ctx.BadRequest()
	// }
	gID := microsoftUser.Email.Address
	if gID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Unable to get Microsoft user ID"))
	}

	switch gc.Purpose {
	case microsoftRegister:
		_, err := database.GetMicrosoftAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Microsoft account is already attached to an account"))
		} else if err != database.ErrMicrosoftAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		gr := &database.MicrosoftRegister{
			MicrosoftEmail: gID,
			TimeCreated:    time.Now(),
		}
		regID, err := database.CreateMicrosoftRegister(ctx, gr)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		grm := &app.MicrosoftRegisterMedia{
			OauthKey:  regID,
			FirstName: microsoftUser.FirstName,
			LastName:  microsoftUser.LastName,
			Email:     microsoftUser.Email.Address,
		}
		return ctx.OK(grm)
	case microsoftLogin:
		account, err := database.GetMicrosoftAccount(ctx, gID)
		if err == database.ErrMicrosoftAccountNotFound {
			return ctx.BadRequest(goa.ErrBadRequest("No account associated with that Microsoft account"))
		} else if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		u, err := database.GetUser(ctx, account.UserID)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *u, nil)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		ctx.ResponseData.Header().Set("X-Session", sesToken)
		ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
		return ctx.OK(database.UserToUser(u))
	case microsoftAttach:
		_, err := database.GetMicrosoftAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Microsoft account is already attached to an account"))
		} else if err != database.ErrMicrosoftAccountNotFound {
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

		account := &database.MicrosoftAccount{
			MicrosoftEmail: gID,
			UserID:         uID,
		}
		err = database.CreateMicrosoftAccount(ctx, account)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	default:
		log.Critical(ctx, "Bad Microsoft receive type")
		return ctx.InternalServerError(goa.ErrInternal("Invalid Microsoft connection type"))
	}

	return ctx.OK(nil)
}

// Register runs the register action.
func (c *MicrosoftController) Register(ctx *app.RegisterMicrosoftContext) error {
	gr, err := database.GetMicrosoftRegister(ctx, ctx.Payload.OauthKey)
	if err == database.ErrMicrosoftRegisterNotFound {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if gr.TimeCreated.Add(microsoftRegisterExpire).Before(time.Now()) {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	}
	_, err = database.GetMicrosoftAccount(ctx, gr.MicrosoftEmail)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This Microsoft account is already attached to an account"))
	} else if err != database.ErrMicrosoftAccountNotFound {
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

	newU := database.UserFromMicrosoftRegisterParams(ctx.Payload)
	uID, err := createUser(ctx, newU, &ctx.Payload.GRecaptchaResponse, ipAddr)
	if err == ErrInvalidRecaptcha {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	account := &database.MicrosoftAccount{
		MicrosoftEmail: gr.MicrosoftEmail,
		UserID:         uID,
	}
	err = database.CreateMicrosoftAccount(ctx, account)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteMicrosoftRegister(ctx, ctx.Payload.OauthKey)
	if err != nil {
		log.Warning(ctx, "Unable to delete Microsoft registration progress, Microsoft Key=%v", ctx.Payload.OauthKey)
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
func (c *MicrosoftController) RegisterURL(ctx *app.RegisterURLMicrosoftContext) error {
	gc := &database.MicrosoftConnection{
		TimeCreated: time.Now(),
		Purpose:     microsoftRegister,
	}
	state, err := database.CreateMicrosoftConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(microsoftConf.AuthCodeURL(state.String())))
}
