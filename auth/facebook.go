package auth

import (
	"encoding/json"
	"fmt"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/log"
	"gigglesearch.org/giggle-auth/utils/secrets"
	"github.com/goadesign/goa"
	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
	"net"
	"time"
)

// FacebookController implements the facebook resource.
type FacebookController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

const (
	facebookRegister = iota
	facebookLogin
	facebookAttach
)

const (
	facebookConnectionExpire = 30 * time.Minute
	facebookRegisterExpire   = time.Hour
)

var facebookConf = &oauth2.Config{
	ClientID:     secrets.FacebookID,
	ClientSecret: secrets.FacebookSecret,
	Scopes: []string{
		"public_profile",
		"email",
	},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://www.facebook.com/v2.12/dialog/oauth",
		TokenURL: "https://graph.facebook.com/v2.12/oauth/access_token",
	},
	RedirectURL: "https://localhost:4000/fb_social",
}

// NewFacebookController creates a facebook controller.
func NewFacebookController(service *goa.Service, jwtSec auth.JWTSecurity, sesCont *SessionController) *FacebookController {
	return &FacebookController{
		Controller:        service.NewController("FacebookController"),
		JWTSecurity:       jwtSec,
		sessionController: sesCont,
	}
}

// AttachToAccount runs the attach-to-account action.
func (c *FacebookController) AttachToAccount(ctx *app.AttachToAccountFacebookContext) error {
	gc := &database.FacebookConnection{
		TimeCreated: time.Now(),
		Purpose:     facebookAttach,
	}
	state, err := database.CreateFacebookConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(facebookConf.AuthCodeURL(state.String())))
}

// DetachFromAccount runs the detach-from-account action.
func (c *FacebookController) DetachFromAccount(ctx *app.DetachFromAccountFacebookContext) error {
	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot detach last login method"))
	}

	gID, err := database.QueryFacebookAccountUser(ctx, uID)
	if err == database.ErrFacebookAccountNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account is not connected to Facebook"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteFacebookAccount(ctx, gID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))
}

// Login runs the login action.
func (c *FacebookController) Login(ctx *app.LoginFacebookContext) error {
	var mt uuid.UUID
	if ctx.Token != nil {
		mt = *ctx.Token
	}
	gc := &database.FacebookConnection{
		TimeCreated: time.Now(),
		Purpose:     facebookLogin,
		MergeToken:  mt,
	}
	state, err := database.CreateFacebookConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(facebookConf.AuthCodeURL(state.String())))
}

type FacebookUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (g *FacebookUser) GetID() string {
	return g.ID
}

// Receive runs the receive action.
func (c *FacebookController) Receive(ctx *app.ReceiveFacebookContext) error {
	gc, err := database.GetFacebookConnection(ctx, ctx.State)
	if err == database.ErrFacebookConnectionNotFound {
		return ctx.BadRequest(goa.ErrBadRequest("Facebook connection must be created with other API methods"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteFacebookConnection(ctx, ctx.State)
	if err != nil {
		log.Warning(ctx, "Unable to delete Facebook connection, state=%s", ctx.State)
	}
	if gc.TimeCreated.Add(facebookConnectionExpire).Before(time.Now()) {
		return ctx.BadRequest(goa.ErrBadRequest("Facebook connection must be created with other API methods"))
	}

	token, err := facebookConf.Exchange(ctx, ctx.Code)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	tokenClient := facebookConf.Client(ctx, token)
	resp, err := tokenClient.Get("https://graph.facebook.com/v2.12/me?fields=id,first_name,last_name,email")
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid Facebook data"))
	}
	facebookUser := &FacebookUser{}
	err = json.NewDecoder(resp.Body).Decode(facebookUser)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	// cli := google.NewClient(tokenClient)
	// googleUser, _, err := cli.Users.Get(ctx, "")
	// if err != nil {
	// 	return ctx.BadRequest()
	// }
	gID := facebookUser.GetID()
	if gID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Unable to get Facebook user ID"))
	}

	switch gc.Purpose {
	case facebookRegister:
		_, err := database.GetFacebookAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Facebook account is already attached to an account"))
		} else if err != database.ErrFacebookAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		gr := &database.FacebookRegister{
			FacebookID:  gID,
			TimeCreated: time.Now(),
		}
		regID, err := database.CreateFacebookRegister(ctx, gr)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		grm := &app.FacebookRegisterMedia{
			OauthKey:  regID,
			Email:     facebookUser.Email,
			FirstName: facebookUser.FirstName,
			LastName:  facebookUser.LastName,
		}
		return ctx.OK(grm)
	case facebookLogin:
		account, err := database.GetFacebookAccount(ctx, gID)
		if err == database.ErrFacebookAccountNotFound {
			return ctx.BadRequest(goa.ErrBadRequest("No account associated with that Facebook account"))
		} else if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		u, err := database.GetUser(ctx, account.UserID)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		var mt *uuid.UUID

		if gc.MergeToken.String() != mt.String() {
			mt = &gc.MergeToken
		}

		sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *u, mt)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		ctx.ResponseData.Header().Set("X-Session", sesToken)
		ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
		return ctx.OK(database.UserToUser(u))
	case facebookAttach:
		_, err := database.GetFacebookAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Facebook account is already attached to an account"))
		} else if err != database.ErrFacebookAccountNotFound {
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

		account := &database.FacebookAccount{
			ID:     gID,
			UserID: uID,
		}
		err = database.CreateFacebookAccount(ctx, account)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	default:
		log.Critical(ctx, "Bad Facebook receive type")
		return ctx.InternalServerError(goa.ErrInternal("Invalid Facebook connection type"))
	}

	return ctx.OK(nil)
}

// Register runs the register action.
func (c *FacebookController) Register(ctx *app.RegisterFacebookContext) error {
	gr, err := database.GetFacebookRegister(ctx, ctx.Payload.OauthKey)
	if err == database.ErrFacebookRegisterNotFound {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if gr.TimeCreated.Add(facebookRegisterExpire).Before(time.Now()) {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	}
	_, err = database.GetFacebookAccount(ctx, gr.FacebookID)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This Facebook account is already attached to an account"))
	} else if err != database.ErrFacebookAccountNotFound {
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

	newU := database.UserFromFacebookRegisterParams(ctx.Payload)
	uID, err := createUser(ctx, newU, &ctx.Payload.GRecaptchaResponse, ipAddr)
	if err == ErrInvalidRecaptcha {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	account := &database.FacebookAccount{
		ID:     gr.FacebookID,
		UserID: uID,
	}
	err = database.CreateFacebookAccount(ctx, account)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteFacebookRegister(ctx, ctx.Payload.OauthKey)
	if err != nil {
		log.Warning(ctx, "Unable to delete Facebook registration progress, Github Key=%v", ctx.Payload.OauthKey)
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
func (c *FacebookController) RegisterURL(ctx *app.RegisterURLFacebookContext) error {
	gc := &database.FacebookConnection{
		TimeCreated: time.Now(),
		Purpose:     facebookRegister,
	}
	state, err := database.CreateFacebookConnection(ctx, gc)
	if err != nil {
		fmt.Println("Error: ", err)
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(facebookConf.AuthCodeURL(state.String())))
}
