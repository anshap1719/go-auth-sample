// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "user": session Resource Client
//
// Command:
// $ goagen
// --design=gigglesearch.org/giggle-auth/auth/design
// --out=$(GOPATH)/src/gigglesearch.org/giggle-auth/auth
// --version=v1.3.1

package client

import (
	"bytes"
	"context"
	"fmt"
	uuid "github.com/goadesign/goa/uuid"
	"net/http"
	"net/url"
)

// CleanLoginTokenSessionPath computes a request path to the clean-login-token action of session.
func CleanLoginTokenSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/clean/token/login")
}

// Cleans old login tokens from the database
func (c *Client) CleanLoginTokenSession(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewCleanLoginTokenSessionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCleanLoginTokenSessionRequest create the request corresponding to the clean-login-token action endpoint of the session resource.
func (c *Client) NewCleanLoginTokenSessionRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// CleanMergeTokenSessionPath computes a request path to the clean-merge-token action of session.
func CleanMergeTokenSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/clean/token/merge")
}

// Cleans old account merge tokens from the database
func (c *Client) CleanMergeTokenSession(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewCleanMergeTokenSessionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCleanMergeTokenSessionRequest create the request corresponding to the clean-merge-token action endpoint of the session resource.
func (c *Client) NewCleanMergeTokenSessionRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// CleanSessionsSessionPath computes a request path to the clean-sessions action of session.
func CleanSessionsSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/clean/sessions")
}

// Deletes all the sessions that have expired
func (c *Client) CleanSessionsSession(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewCleanSessionsSessionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCleanSessionsSessionRequest create the request corresponding to the clean-sessions action endpoint of the session resource.
func (c *Client) NewCleanSessionsSessionRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// GetSessionsSessionPath computes a request path to the get-sessions action of session.
func GetSessionsSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/sessions")
}

// Gets all of the sessions that are associated with the currently logged in user
func (c *Client) GetSessionsSession(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetSessionsSessionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetSessionsSessionRequest create the request corresponding to the get-sessions action endpoint of the session resource.
func (c *Client) NewGetSessionsSessionRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// LogoutSessionPath computes a request path to the logout action of session.
func LogoutSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/logout")
}

// Takes a user's auth token, and logs-out the session associated with it
func (c *Client) LogoutSession(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewLogoutSessionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewLogoutSessionRequest create the request corresponding to the logout action endpoint of the session resource.
func (c *Client) NewLogoutSessionRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// LogoutOtherSessionPath computes a request path to the logout-other action of session.
func LogoutOtherSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/logout/all")
}

// Logout all sessions for the current user except their current session
func (c *Client) LogoutOtherSession(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewLogoutOtherSessionRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewLogoutOtherSessionRequest create the request corresponding to the logout-other action endpoint of the session resource.
func (c *Client) NewLogoutOtherSessionRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// LogoutSpecificSessionPath computes a request path to the logout-specific action of session.
func LogoutSpecificSessionPath(session string) string {
	param0 := session

	return fmt.Sprintf("/api/v1/user/auth/logout/%s-id", param0)
}

// Logout of a specific session
func (c *Client) LogoutSpecificSession(ctx context.Context, path string, sessionID string) (*http.Response, error) {
	req, err := c.NewLogoutSpecificSessionRequest(ctx, path, sessionID)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewLogoutSpecificSessionRequest create the request corresponding to the logout-specific action endpoint of the session resource.
func (c *Client) NewLogoutSpecificSessionRequest(ctx context.Context, path string, sessionID string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("session-id", sessionID)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// RedeemTokenSessionPayload is the session redeemToken action payload.
type RedeemTokenSessionPayload struct {
	// The token to redeem
	Token uuid.UUID `form:"token" json:"token" yaml:"token" xml:"token"`
}

// RedeemTokenSessionPath computes a request path to the redeemToken action of session.
func RedeemTokenSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/token")
}

// Redeems a login token for credentials
func (c *Client) RedeemTokenSession(ctx context.Context, path string, payload *RedeemTokenSessionPayload, contentType string) (*http.Response, error) {
	req, err := c.NewRedeemTokenSessionRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRedeemTokenSessionRequest create the request corresponding to the redeemToken action endpoint of the session resource.
func (c *Client) NewRedeemTokenSessionRequest(ctx context.Context, path string, payload *RedeemTokenSessionPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// RefreshSessionPath computes a request path to the refresh action of session.
func RefreshSessionPath() string {

	return fmt.Sprintf("/api/v1/user/auth/session")
}

// Take a user's session token and refresh it, also returns a new authentication token
func (c *Client) RefreshSession(ctx context.Context, path string, xSession string) (*http.Response, error) {
	req, err := c.NewRefreshSessionRequest(ctx, path, xSession)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRefreshSessionRequest create the request corresponding to the refresh action endpoint of the session resource.
func (c *Client) NewRefreshSessionRequest(ctx context.Context, path string, xSession string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, err
	}
	header := req.Header

	header.Set("X-Session", xSession)

	return req, nil
}