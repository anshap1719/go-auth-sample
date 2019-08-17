package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/auth/models"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/log"
	"gigglesearch.org/giggle-auth/utils/secrets"
	"github.com/goadesign/goa"
	"github.com/gofrs/uuid"
	"github.com/mssola/user_agent"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	SessionTime = 7 * 24 * time.Hour // 1 week
	TokenTime   = 10 * time.Minute   // 10 minutes
)

// SessionController implements the session resource.
type SessionController struct {
	*goa.Controller
	auth.JWTSecurity
}

// NewSessionController creates a session controller.
func NewSessionController(service *goa.Service, jwtSec auth.JWTSecurity) *SessionController {
	return &SessionController{
		Controller:  service.NewController("SessionController"),
		JWTSecurity: jwtSec,
	}
}

// CleanLoginToken runs the clean-login-token action.
func (c *SessionController) CleanLoginToken(ctx *app.CleanLoginTokenSessionContext) error {
	// SessionController_CleanLoginToken: start_implement

	if c.IsAdmin(ctx.Request) {
		tokens, err := database.QueryLoginTokenOld(ctx, time.Now())
		if err != nil {
			return ctx.OK([]byte(""))
		}

		for start := 0; start < len(tokens); start += 500 {
			end := start + 500
			if end > len(tokens) {
				end = len(tokens)
			}
			_ = database.DeleteLoginTokenMulti(ctx, tokens[start:end])
		}

		return ctx.OK([]byte(""))
	} else {
		return ctx.Forbidden(goa.ErrUnauthorized("user is not admin"))
	}

	// SessionController_CleanLoginToken: end_implement
}

// CleanMergeToken runs the clean-merge-token action.
func (c *SessionController) CleanMergeToken(ctx *app.CleanMergeTokenSessionContext) error {
	// SessionController_CleanMergeToken: start_implement
	if c.IsAdmin(ctx.Request) {
		tokens, err := database.QueryMergeTokenOld(ctx, time.Now())
		if err != nil {
			return ctx.OK([]byte(""))
		}

		for start := 0; start < len(tokens); start += 500 {
			end := start + 500
			if end > len(tokens) {
				end = len(tokens)
			}
			_ = database.DeleteMergeTokenMulti(ctx, tokens[start:end])
		}

		return ctx.OK([]byte(""))
	} else {
		return ctx.Forbidden(goa.ErrUnauthorized("user is not admin"))
	}

	// SessionController_CleanMergeToken: end_implement
}

// CleanSessions runs the clean-sessions action.
func (c *SessionController) CleanSessions(ctx *app.CleanSessionsSessionContext) error {
	// SessionController_CleanSessions: start_implement

	if c.IsAdmin(ctx.Request) {
		sessionIds, err := database.QuerySessionOld(ctx, time.Now().Add(-SessionTime))
		if err != nil {
			return ctx.OK([]byte(""))
		}

		for start := 0; start < len(sessionIds); start += 500 {
			end := start + 500
			if end > len(sessionIds) {
				end = len(sessionIds)
			}
			_ = database.DeleteSessionMulti(ctx, sessionIds[start:end])
		}

		return ctx.OK([]byte(""))
	} else {
		return ctx.Forbidden(goa.ErrUnauthorized("user is not admin"))
	}

	// SessionController_CleanSessions: end_implement
}

// GetSessions runs the get-sessions action.
func (c *SessionController) GetSessions(ctx *app.GetSessionsSessionContext) error {
	// SessionController_GetSessions: start_implement

	userID := c.GetUserID(ctx.Request)
	sesID := c.GetSessionFromAuth(ctx.Request)
	sessions, err := database.QuerySessionFromAccount(ctx, userID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	capacity := len(sessions) - 1
	if capacity < 0 {
		capacity = 0
	}
	res := &app.AllSessions{
		OtherSessions: make(app.SessionCollection, 0, capacity),
	}
	for _, v := range sessions {
		s := database.SessionToSession(v)
		if v.Coordinates != "" {
			mapURL, err := getMapURL(v.Coordinates)
			if err != nil {
				return ctx.InternalServerError(goa.ErrInternal(err))
			}
			s.MapURL = mapURL
		}
		if v.ID.Hex() == sesID {
			res.CurrentSession = s
		} else {
			res.OtherSessions = append(res.OtherSessions, s)
		}
	}
	return ctx.OK(res)

	// SessionController_GetSessions: end_implement
}

// Logout runs the logout action.
func (c *SessionController) Logout(ctx *app.LogoutSessionContext) error {
	// SessionController_Logout: start_implement

	sesID := c.GetSessionFromAuth(ctx.Request)

	if err := database.DeleteSession(ctx, sesID); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))

	// SessionController_Logout: end_implement
}

// LogoutOther runs the logout-other action.
func (c *SessionController) LogoutOther(ctx *app.LogoutOtherSessionContext) error {
	uID := c.GetUserID(ctx.Request)
	sesID := c.GetSessionFromAuth(ctx.Request)

	err := logoutAllSessionsBut(ctx, uID, sesID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))
}

// LogoutSpecific runs the logout-specific action.
func (c *SessionController) LogoutSpecific(ctx *app.LogoutSpecificSessionContext) error {
	uID := c.GetUserID(ctx.Request)
	sesID := ctx.SessionID
	if sesID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Session ID must be provided"))
	}

	s, err := database.GetSession(ctx, sesID)
	if err == database.ErrSessionNotFound {
		return ctx.NotFound(goa.ErrNotFound("No session with the given ID found"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if s.UserID != uID {
		return ctx.NotFound(goa.ErrNotFound("No session with the given ID found"))
	}

	err = database.DeleteSession(ctx, sesID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))
}

// RedeemToken runs the redeemToken action.
func (c *SessionController) RedeemToken(ctx *app.RedeemTokenSessionContext) error {
	// SessionController_RedeemToken: start_implement

	t, err := database.GetLoginToken(ctx, ctx.Payload.Token)
	if err == database.ErrLoginTokenNotFound {
		return ctx.Forbidden(errTokenMismatch("Token does not exist"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if t.TimeExpire.Before(time.Now()) {
		return ctx.Forbidden(errTokenMismatch("Token does not exist"))
	}

	user, err := database.GetUser(ctx, t.UserID)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	sesToken, authToken, err := c.loginUser(ctx, ctx.Request, *user, nil)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = database.DeleteLoginToken(ctx, t.Token)
	if err != nil {
		log.Warning(ctx, "Unable to delete login token %s", t.Token)
	}

	ctx.ResponseData.Header().Set("X-Session", sesToken)
	ctx.ResponseData.Header().Set("Authorization", authToken)
	return ctx.Created()

	// SessionController_RedeemToken: end_implement
}

// Refresh runs the refresh action.
func (c *SessionController) Refresh(ctx *app.RefreshSessionContext) error {
	// SessionController_Refresh: start_implement

	sesID := c.GetSessionCode(ctx.Request)
	if sesID == "" {
		return ctx.BadRequest(goa.ErrBadRequest("Invalid session ID"))
	}
	s, err := database.GetSession(ctx, sesID)
	if err == database.ErrSessionNotFound {
		return ctx.Unauthorized(goa.ErrUnauthorized("Session not found"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	if !s.LastUsed.Add(SessionTime).After(time.Now()) {
		return ctx.Unauthorized(goa.ErrUnauthorized("Session not found"))
	}

	updatedSession := c.createSession(ctx.Request, s.UserID, s.IsAdmin, s.IsPluginAuthor, s.IsEventAuthor)
	updatedSession.ID = s.ID
	err = database.UpdateSession(ctx, updatedSession)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	sessToken, err := c.SignSessionToken(SessionTime, s.ID.Hex())
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	authToken, err := c.SignAuthToken(TokenTime, s.ID.Hex(), s.UserID, s.IsAdmin, s.IsPluginAuthor, s.IsEventAuthor)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	ctx.ResponseData.Header().Set("X-Session", sessToken)
	ctx.ResponseData.Header().Set("Authorization", "Bearer "+authToken)
	return ctx.OK([]byte(""))

	// SessionController_Refresh: end_implement
}

func (c *SessionController) createSession(req *http.Request, userID string, isAdmin, isPluginAuthor, isEventAuthor bool) *database.Session {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ip = req.RemoteAddr
	}

	ua := user_agent.New(req.UserAgent())
	os := ua.OSInfo()
	if os.Name == "OS" {
		os.Name = ua.Platform()
	}
	browse, vers := ua.Browser()

	//location, coord := getIPLocation(req.Header)

	newSession := &database.Session{
		UserID:         userID,
		LastUsed:       time.Now(),
		IP:             ip,
		Os:             os.Name + " " + os.Version,
		Browser:        browse + " " + vers,
		Location:       "",
		Coordinates:    "",
		IsMobile:       ua.Mobile(),
		IsAdmin:        isAdmin,
		IsPluginAuthor: isPluginAuthor,
		IsEventAuthor:  isEventAuthor,
	}
	return newSession
}

func (c *SessionController) loginUser(ctx context.Context, req *http.Request, user models.User, mergeToken *uuid.UUID) (sessionToken string, authToken string, err error) {
	newSession := c.createSession(req, user.ID.Hex(), user.IsAdmin, user.IsPluginAuthor, user.IsEventAuthor)
	sesID, err := database.CreateSession(ctx, newSession)
	if err != nil {
		return "", "", err
	}

	sessionToken, err = c.SignSessionToken(SessionTime, sesID)
	if err != nil {
		return "", "", err
	}

	authToken, err = c.SignAuthToken(TokenTime, sesID, user.ID.Hex(), user.IsAdmin, user.IsPluginAuthor, user.IsEventAuthor)
	if err != nil {
		return "", "", err
	}

	return sessionToken, authToken, nil
}

func logoutAllSessionsBut(ctx context.Context, userID, sessionID string) error {
	sesIds, err := database.QuerySessionIds(ctx, userID)
	if err != nil {
		return err
	}

	for i, v := range sesIds {
		if v == sessionID {
			sesIds = append(sesIds[:i], sesIds[i+1:]...)
			break
		}
	}

	err = database.DeleteSessionMulti(ctx, sesIds)
	if err != nil {
		return err
	}
	return nil
}

//func getIPLocation(headerInfo http.Header) (string, string) {
//	country := headerInfo.Get("X-Appengine-Country")
//	if country == "ZZ" || country == "" {
//		return "Unknown", ""
//	}
//
//	region := headerInfo.Get("X-AppEngine-Region")
//
//	city := strings.Title(headerInfo.Get("X-AppEngine-City"))
//
//	latlong := headerInfo.Get("X-AppEngine-CityLatLong")
//
//	currentLoc, ok := locMap[country]
//	if !ok {
//		return "Unknown", ""
//	}
//
//	currentRegion, ok := currentLoc.Region[region]
//	if !ok {
//		return city + ", " + currentLoc.Name, latlong
//	}
//
//	return city + ", " + currentRegion + ", " + currentLoc.Name, latlong
//}
//
func getMapURL(coordinates string) (string, error) {
	u, err := url.Parse("https://maps.googleapis.com/maps/api/staticmap?size=500x500&markers=" + coordinates + "&format=jpg&zoom=7&key=" + secrets.GoogleMapsKey)
	if err != nil {
		return "", err
	}

	sign, err := base64.URLEncoding.DecodeString(secrets.GoogleMapsSigning)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha1.New, sign)
	_, err = io.WriteString(h, u.RequestURI())
	if err != nil {
		return "", err
	}

	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	u.RawQuery += "&signature=" + signature
	return u.String(), nil
}
