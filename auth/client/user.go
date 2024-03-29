// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "user": user Resource Client
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
	"net/http"
	"net/url"
	"strconv"
)

// AddPluginUserPath computes a request path to the add-plugin action of user.
func AddPluginUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/plugins")
}

// Add a new plugin to user's account
func (c *Client) AddPluginUser(ctx context.Context, path string, payload *UserPlugin, contentType string) (*http.Response, error) {
	req, err := c.NewAddPluginUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewAddPluginUserRequest create the request corresponding to the add-plugin action endpoint of the user resource.
func (c *Client) NewAddPluginUserRequest(ctx context.Context, path string, payload *UserPlugin, contentType string) (*http.Request, error) {
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

// DeactivateUserPath computes a request path to the deactivate action of user.
func DeactivateUserPath() string {

	return fmt.Sprintf("/api/v1/user/user")
}

// Disable a user's account
func (c *Client) DeactivateUser(ctx context.Context, path string, admin *bool, id *string) (*http.Response, error) {
	req, err := c.NewDeactivateUserRequest(ctx, path, admin, id)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewDeactivateUserRequest create the request corresponding to the deactivate action endpoint of the user resource.
func (c *Client) NewDeactivateUserRequest(ctx context.Context, path string, admin *bool, id *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if admin != nil {
		tmp100 := strconv.FormatBool(*admin)
		values.Set("admin", tmp100)
	}
	if id != nil {
		values.Set("id", *id)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("DELETE", u.String(), nil)
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

// GetAllUsersUserPath computes a request path to the get-all-users action of user.
func GetAllUsersUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/all")
}

// Get all users
func (c *Client) GetAllUsersUser(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetAllUsersUserRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetAllUsersUserRequest create the request corresponding to the get-all-users action endpoint of the user resource.
func (c *Client) NewGetAllUsersUserRequest(ctx context.Context, path string) (*http.Request, error) {
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

// GetByEmailUserPath computes a request path to the get-by-email action of user.
func GetByEmailUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/email")
}

// Get a user by their email. Only callable by admins
func (c *Client) GetByEmailUser(ctx context.Context, path string, email string) (*http.Response, error) {
	req, err := c.NewGetByEmailUserRequest(ctx, path, email)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetByEmailUserRequest create the request corresponding to the get-by-email action endpoint of the user resource.
func (c *Client) NewGetByEmailUserRequest(ctx context.Context, path string, email string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	values.Set("email", email)
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

// GetManyUserPath computes a request path to the get-many action of user.
func GetManyUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/multi")
}

// Get many users by their ID
func (c *Client) GetManyUser(ctx context.Context, path string, id []string) (*http.Response, error) {
	req, err := c.NewGetManyUserRequest(ctx, path, id)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetManyUserRequest create the request corresponding to the get-many action endpoint of the user resource.
func (c *Client) NewGetManyUserRequest(ctx context.Context, path string, id []string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	for _, p := range id {
		tmp101 := p
		values.Add("id", tmp101)
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

// GetAuthsUserPath computes a request path to the getAuths action of user.
func GetAuthsUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/authstat")
}

// Returns whether Oauth is attached or not
func (c *Client) GetAuthsUser(ctx context.Context, path string, userID *string) (*http.Response, error) {
	req, err := c.NewGetAuthsUserRequest(ctx, path, userID)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetAuthsUserRequest create the request corresponding to the getAuths action endpoint of the user resource.
func (c *Client) NewGetAuthsUserRequest(ctx context.Context, path string, userID *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if userID != nil {
		values.Set("userID", *userID)
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

// ResendVerifyEmailUserPath computes a request path to the resend-verify-email action of user.
func ResendVerifyEmailUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/resend-verify")
}

// Resends a verify email for the current user, also invalidates the link on the previously send email verification
func (c *Client) ResendVerifyEmailUser(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewResendVerifyEmailUserRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewResendVerifyEmailUserRequest create the request corresponding to the resend-verify-email action endpoint of the user resource.
func (c *Client) NewResendVerifyEmailUserRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), nil)
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

// RetrieveUserPath computes a request path to the retrieve action of user.
func RetrieveUserPath() string {

	return fmt.Sprintf("/api/v1/user/user")
}

// Get user by ID
func (c *Client) RetrieveUser(ctx context.Context, path string, userID *string) (*http.Response, error) {
	req, err := c.NewRetrieveUserRequest(ctx, path, userID)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewRetrieveUserRequest create the request corresponding to the retrieve action endpoint of the user resource.
func (c *Client) NewRetrieveUserRequest(ctx context.Context, path string, userID *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if userID != nil {
		values.Set("user-id", *userID)
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

// UpdateUserPath computes a request path to the update action of user.
func UpdateUserPath() string {

	return fmt.Sprintf("/api/v1/user/user")
}

// Update a user
func (c *Client) UpdateUser(ctx context.Context, path string, payload *UserParams, contentType string) (*http.Response, error) {
	req, err := c.NewUpdateUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUpdateUserRequest create the request corresponding to the update action endpoint of the user resource.
func (c *Client) NewUpdateUserRequest(ctx context.Context, path string, payload *UserParams, contentType string) (*http.Request, error) {
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
	req, err := http.NewRequest("PATCH", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// UpdateAdminUserPath computes a request path to the update-admin action of user.
func UpdateAdminUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/update-user")
}

// Update a user from admin dashboard
func (c *Client) UpdateAdminUser(ctx context.Context, path string, payload *UserParamsAdmin, uid *string, contentType string) (*http.Response, error) {
	req, err := c.NewUpdateAdminUserRequest(ctx, path, payload, uid, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUpdateAdminUserRequest create the request corresponding to the update-admin action endpoint of the user resource.
func (c *Client) NewUpdateAdminUserRequest(ctx context.Context, path string, payload *UserParamsAdmin, uid *string, contentType string) (*http.Request, error) {
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
	values := u.Query()
	if uid != nil {
		values.Set("uid", *uid)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("PATCH", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	if c.JWTSigner != nil {
		if err := c.JWTSigner.Sign(req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

// UpdatePluginPermissionsUserPath computes a request path to the update-plugin-permissions action of user.
func UpdatePluginPermissionsUserPath() string {

	return fmt.Sprintf("/api/v1/user/user/plugins")
}

// Update plugin permissions
func (c *Client) UpdatePluginPermissionsUser(ctx context.Context, path string, payload *UserPlugin, contentType string) (*http.Response, error) {
	req, err := c.NewUpdatePluginPermissionsUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUpdatePluginPermissionsUserRequest create the request corresponding to the update-plugin-permissions action endpoint of the user resource.
func (c *Client) NewUpdatePluginPermissionsUserRequest(ctx context.Context, path string, payload *UserPlugin, contentType string) (*http.Request, error) {
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
	req, err := http.NewRequest("PUT", u.String(), &body)
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

// ValidateEmailUserPath computes a request path to the validate-email action of user.
func ValidateEmailUserPath(validateID string) string {
	param0 := validateID

	return fmt.Sprintf("/verifyemail/%s", param0)
}

// Validates an email address, designed to be called by users directly in their browser
func (c *Client) ValidateEmailUser(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewValidateEmailUserRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewValidateEmailUserRequest create the request corresponding to the validate-email action endpoint of the user resource.
func (c *Client) NewValidateEmailUserRequest(ctx context.Context, path string) (*http.Request, error) {
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
