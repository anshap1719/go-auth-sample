// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "user": Application Security
//
// Command:
// $ goagen
// --design=gigglesearch.org/giggle-auth/auth/design
// --out=$(GOPATH)/src/gigglesearch.org/giggle-auth/auth
// --version=v1.3.1

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

type (
	// Private type used to store auth handler info in request context
	authMiddlewareKey string
)

// UseJWTMiddleware mounts the jwt auth middleware onto the service.
func UseJWTMiddleware(service *goa.Service, middleware goa.Middleware) {
	service.Context = context.WithValue(service.Context, authMiddlewareKey("jwt"), middleware)
}

// NewJWTSecurity creates a jwt security definition.
func NewJWTSecurity() *goa.JWTSecurity {
	def := goa.JWTSecurity{
		In:       goa.LocHeader,
		Name:     "Authorization",
		TokenURL: "http://localhost:3000/api/auth/login",
	}
	return &def
}

// UseKeyMiddleware mounts the key auth middleware onto the service.
func UseKeyMiddleware(service *goa.Service, middleware goa.Middleware) {
	service.Context = context.WithValue(service.Context, authMiddlewareKey("key"), middleware)
}

// NewKeySecurity creates a key security definition.
func NewKeySecurity() *goa.APIKeySecurity {
	def := goa.APIKeySecurity{
		In:   goa.LocHeader,
		Name: "API-Key",
	}
	def.Description = "API Key for users API"
	return &def
}

// handleSecurity creates a handler that runs the auth middleware for the security scheme.
func handleSecurity(schemeName string, h goa.Handler, scopes ...string) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		scheme := ctx.Value(authMiddlewareKey(schemeName))
		am, ok := scheme.(goa.Middleware)
		if !ok {
			return goa.NoAuthMiddleware(schemeName)
		}
		ctx = goa.WithRequiredScopes(ctx, scopes)
		return am(h)(ctx, rw, req)
	}
}
