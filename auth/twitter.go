package auth

import (
	"encoding/json"
	"fmt"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/auth/models"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/log"
	"gigglesearch.org/giggle-auth/utils/secrets"
	"github.com/globalsign/mgo"
	"github.com/goadesign/goa"
	"github.com/gofrs/uuid"
	"github.com/mrjones/oauth"
	"net"
	"strings"
	"time"
)

const (
	twitterRegister = iota
	twitterLogin
	twitterAttach
)

const (
	twitterConnectionExpire = 30 * time.Minute
	twitterRegisterExpire   = time.Hour
)

var twitterConf = oauth.NewConsumer(
	secrets.TwitterKey,
	secrets.TwitterSecret,
	oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	})

// TwitterController implements the twitter resource.
type TwitterController struct {
	*goa.Controller
	auth.JWTSecurity
	sessionController *SessionController
}

// NewTwitterController creates a twitter controller.
func NewTwitterController(service *goa.Service, jwtSec auth.JWTSecurity, session *SessionController) *TwitterController {
	twitterConf.Debug(true)

	return &TwitterController{
		Controller:        service.NewController("TwitterController"),
		JWTSecurity:       jwtSec,
		sessionController: session,
	}
}

// AttachToAccount runs the attach-to-account action.
func (c *TwitterController) AttachToAccount(ctx *app.AttachToAccountTwitterContext) error {
	// TwitterController_AttachToAccount: start_implement

	gc := &database.TwitterConnection{
		TimeCreated: time.Now(),
		Purpose:     twitterAttach,
	}
	state, err := database.CreateTwitterConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	tokenUrl := "https://localhost:4000/tw_social?state=" + state.String()

	token, requestUrl, err := twitterConf.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if err := database.CreateTwitterToken(token.Token, *token); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(requestUrl))

	// TwitterController_AttachToAccount: end_implement
}

// DetachFromAccount runs the detach-from-account action.
func (c *TwitterController) DetachFromAccount(ctx *app.DetachFromAccountTwitterContext) error {
	// TwitterController_DetachFromAccount: start_implement

	uID := c.GetUserID(ctx.Request)

	if getNumLoginMethods(ctx, uID) <= 1 {
		return ctx.Forbidden(errMustBeAbleToLogin("Cannot detach last login method"))
	}

	gID, err := database.QueryTwitterAccountUser(ctx, uID)
	if err == database.ErrTwitterAccountNotFound {
		return ctx.NotFound(goa.ErrNotFound("User account is not connected to Twitter"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteTwitterAccount(ctx, gID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	return ctx.OK([]byte(""))

	// TwitterController_DetachFromAccount: end_implement
}

// Login runs the login action.
func (c *TwitterController) Login(ctx *app.LoginTwitterContext) error {
	// TwitterController_Login: start_implement

	var mt uuid.UUID
	if ctx.Token != nil {
		mt = *ctx.Token
	}
	gc := &database.TwitterConnection{
		TimeCreated: time.Now(),
		Purpose:     twitterLogin,
		MergeToken:  mt,
	}
	state, err := database.CreateTwitterConnection(ctx, gc)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	tokenUrl := "https://localhost:4000/tw_social?state=" + state.String()

	token, requestUrl, err := twitterConf.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if err := database.CreateTwitterToken(token.Token, *token); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(requestUrl))

	// TwitterController_Login: end_implement
}

// Receive runs the receive action.
func (c *TwitterController) Receive(ctx *app.ReceiveTwitterContext) error {
	// TwitterController_Receive: start_implement

	state, err := uuid.FromString(ctx.State)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest("State UUID is invalid"))
	}

	gc, err := database.GetTwitterConnection(ctx, state)
	if err == database.ErrTwitterConnectionNotFound {
		return ctx.BadRequest(goa.ErrBadRequest("Twitter connection must be created with other API methods"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteTwitterConnection(ctx, state)
	if err != nil {
		log.Warning(ctx, "Unable to delete Twitter connection, state=%s", ctx.State)
	}
	if gc.TimeCreated.Add(twitterConnectionExpire).Before(time.Now()) {
		return ctx.BadRequest(goa.ErrBadRequest("Twitter connection must be created with other API methods"))
	}

	code := ctx.OauthVerifier
	key := ctx.OauthToken

	token, err := database.GetTwitterToken(key)
	if err == mgo.ErrNotFound {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = database.DeleteTwitterToken(key)
	if err != nil {
		log.Warning(ctx, "Unable to delete Twitter connection, key=%s", key)
	}

	accessToken, err := twitterConf.AuthorizeToken(token, code)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	client, err := twitterConf.MakeHttpClient(accessToken)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	response, err := client.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json?include_entities=false&skip_status=true&include_email=true")
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	defer response.Body.Close()

	var twitterUser models.TwitterResponse

	if err := json.NewDecoder(response.Body).Decode(&twitterUser); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	gID := string(twitterUser.ID)

	switch gc.Purpose {
	case twitterRegister:
		_, err := database.GetTwitterAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Twitter account is already attached to an account"))
		} else if err != database.ErrTwitterAccountNotFound {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		gr := &database.TwitterRegister{
			TwitterID:   gID,
			TimeCreated: time.Now(),
		}
		regID, err := database.CreateTwitterRegister(ctx, gr)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		grm := &app.TwitterRegisterMedia{
			OauthKey:  regID,
			Email:     twitterUser.Email,
			FirstName: strings.Split(twitterUser.Name, " ")[0],
			LastName:  strings.Split(twitterUser.Name, " ")[1],
		}
		return ctx.OK(grm)
	case twitterLogin:
		account, err := database.GetTwitterAccount(ctx, gID)
		if err == database.ErrTwitterAccountNotFound {
			return ctx.BadRequest(goa.ErrBadRequest("No account associated with that Twitter account"))
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
	case twitterAttach:
		_, err := database.GetTwitterAccount(ctx, gID)
		if err == nil {
			return ctx.BadRequest(goa.ErrBadRequest("This Twitter account is already attached to an account"))
		} else if err != database.ErrTwitterAccountNotFound {
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

		account := &database.TwitterAccount{
			ID:     gID,
			UserID: uID,
		}
		err = database.CreateTwitterAccount(ctx, account)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	default:
		log.Critical(ctx, "Bad Twitter receive type")
		return ctx.InternalServerError(goa.ErrInternal("Invalid Twitter connection type"))
	}

	return ctx.OK(nil)

	// TwitterController_Receive: end_implement
}

// Register runs the register action.
func (c *TwitterController) Register(ctx *app.RegisterTwitterContext) error {
	// TwitterController_Register: start_implement

	gr, err := database.GetTwitterRegister(ctx, ctx.Payload.OauthKey)
	if err == database.ErrTwitterRegisterNotFound {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if gr.TimeCreated.Add(twitterRegisterExpire).Before(time.Now()) {
		return ctx.NotFound(goa.ErrNotFound("Invalid registration key"))
	}
	_, err = database.GetTwitterAccount(ctx, gr.TwitterID)
	if err == nil {
		return ctx.Forbidden(errAlreadyExists("This Twitter account is already attached to an account"))
	} else if err != database.ErrTwitterAccountNotFound {
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

	newU := database.UserFromTwitterRegisterParams(ctx.Payload)
	uID, err := createUser(ctx, newU, &ctx.Payload.GRecaptchaResponse, ipAddr)
	if err == ErrInvalidRecaptcha {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	account := &database.TwitterAccount{
		ID:     gr.TwitterID,
		UserID: uID,
	}
	err = database.CreateTwitterAccount(ctx, account)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	err = database.DeleteTwitterRegister(ctx, ctx.Payload.OauthKey)
	if err != nil {
		log.Warning(ctx, "Unable to delete Twitter registration progress, Github Key=%v", ctx.Payload.OauthKey)
	}

	sesToken, authToken, err := c.sessionController.loginUser(ctx, ctx.Request, *newU, nil)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
	return ctx.OK(database.UserToUser(newU))

	// TwitterController_Register: end_implement
}

// RegisterURL runs the register-url action.
func (c *TwitterController) RegisterURL(ctx *app.RegisterURLTwitterContext) error {
	// TwitterController_RegisterURL: start_implement

	gc := &database.TwitterConnection{
		TimeCreated: time.Now(),
		Purpose:     twitterRegister,
	}
	state, err := database.CreateTwitterConnection(ctx, gc)
	if err != nil {
		fmt.Println("Error: ", err)
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	tokenUrl := "https://localhost:4000/tw_social?state=" + state.String()

	token, requestUrl, err := twitterConf.GetRequestTokenAndUrl(tokenUrl)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if err := database.CreateTwitterToken(token.Token, *token); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(requestUrl))

	// TwitterController_RegisterURL: end_implement
}
