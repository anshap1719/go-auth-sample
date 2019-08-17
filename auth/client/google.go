// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "user": google Resource Client
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

// AttachToAccountGooglePath computes a request path to the attach-to-account action of google.
func AttachToAccountGooglePath() string {

	return fmt.Sprintf("/api/v1/user/auth/google/attach")
}

// Attaches a Google account to an existing user account, returns the URL the browser should be redirected to
func (c *Client) AttachToAccountGoogle(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewAttachToAccountGoogleRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAttachToAccountGoogleRequest create the request corresponding to the attach-to-account action endpoint of the google resource.
func (c *Client) NewAttachToAccountGoogleRequest(ctx context.Context, path string) (*http.Request, error) {
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

// DetachFromAccountGooglePath computes a request path to the detach-from-account action of google.
func DetachFromAccountGooglePath() string {

	return fmt.Sprintf("/api/v1/user/auth/google/detach")
}

// Detaches a Google account from an existing user account.
func (c *Client) DetachFromAccountGoogle(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewDetachFromAccountGoogleRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDetachFromAccountGoogleRequest create the request corresponding to the detach-from-account action endpoint of the google resource.
func (c *Client) NewDetachFromAccountGoogleRequest(ctx context.Context, path string) (*http.Request, error) {
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

// LoginGooglePath computes a request path to the login action of google.
func LoginGooglePath() string {

	return fmt.Sprintf("/api/v1/user/auth/google/login")
}

// Gets the URL the front-end should redirect the browser to in order to be authenticated with Google, to be logged in
func (c *Client) LoginGoogle(ctx context.Context, path string, token *uuid.UUID) (*http.Response, error) {
	req, err := c.NewLoginGoogleRequest(ctx, path, token)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewLoginGoogleRequest create the request corresponding to the login action endpoint of the google resource.
func (c *Client) NewLoginGoogleRequest(ctx context.Context, path string, token *uuid.UUID) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if token != nil {
		tmp92 := token.String()
		values.Set("token", tmp92)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.KeySigner != nil {
		if err := c.KeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// ReceiveGooglePath computes a request path to the receive action of google.
func ReceiveGooglePath() string {

	return fmt.Sprintf("/api/v1/user/auth/google/receive")
}

// The endpoint that Google redirects the browser to after the user has authenticated
func (c *Client) ReceiveGoogle(ctx context.Context, path string, code string, state uuid.UUID) (*http.Response, error) {
	req, err := c.NewReceiveGoogleRequest(ctx, path, code, state)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewReceiveGoogleRequest create the request corresponding to the receive action endpoint of the google resource.
func (c *Client) NewReceiveGoogleRequest(ctx context.Context, path string, code string, state uuid.UUID) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("code", code)
	tmp93 := state.String()
	values.Set("state", tmp93)
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.KeySigner != nil {
		if err := c.KeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// RegisterGooglePath computes a request path to the register action of google.
func RegisterGooglePath() string {

	return fmt.Sprintf("/api/v1/user/auth/google/register")
}

// Registers a new account with the system, with Google as the login system
func (c *Client) RegisterGoogle(ctx context.Context, path string, payload *GoogleRegisterParams, contentType string) (*http.Response, error) {
	req, err := c.NewRegisterGoogleRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRegisterGoogleRequest create the request corresponding to the register action endpoint of the google resource.
func (c *Client) NewRegisterGoogleRequest(ctx context.Context, path string, payload *GoogleRegisterParams, contentType string) (*http.Request, error) {
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
	if c.KeySigner != nil {
		if err := c.KeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// RegisterURLGooglePath computes a request path to the register-url action of google.
func RegisterURLGooglePath() string {

	return fmt.Sprintf("/api/v1/user/auth/google/register-start")
}

// Gets the URL the front-end should redirect the browser to in order to be authenticated with Google, and then register
func (c *Client) RegisterURLGoogle(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRegisterURLGoogleRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRegisterURLGoogleRequest create the request corresponding to the register-url action endpoint of the google resource.
func (c *Client) NewRegisterURLGoogleRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.KeySigner != nil {
		if err := c.KeySigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}