//go:generate goagen bootstrap -d gigglesearch.org/giggle-auth/auth/design

package auth

import (
	"context"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/models"
	"gigglesearch.org/giggle-auth/utils/auth"
	"gigglesearch.org/giggle-auth/utils/database"
	"gigglesearch.org/giggle-auth/utils/log"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	log2 "log"
	"net/http"
)

var errAlreadyExists = goa.NewErrorClass("already-exists", http.StatusForbidden)
var errMustBeAbleToLogin = goa.NewErrorClass("cannot-remove-last-login-option", http.StatusForbidden)
var errDisposable = goa.NewErrorClass("disposable-email", http.StatusForbidden)
var errBadReset = goa.NewErrorClass("bad-reset", http.StatusForbidden)
var errTokenMismatch = goa.NewErrorClass("token-mismatch", http.StatusForbidden)

func main() {
	// Create service
	service := goa.New("user")

	initDB()

	// Setup logging
	logger := log.NewLogger(nil)
	service.WithLogger(logger)

	jwtSec, err := auth.NewJWTSecurity()
	if err != nil {
		log.Critical(context.WithValue(context.Background(), "%v", err), "%v", err)
	}

	jwtMiddle, err := auth.NewJWTMiddleware(app.NewJWTSecurity())
	if err != nil {
		panic(err)
	}

	keyMiddle, err := auth.NewKeyMiddleware()
	if err != nil {
		panic(err)
	}

	// Mount middleware
	service.Use(log.LogContextMiddleware)
	service.Use(log.LogInternalError)
	service.Use(middleware.RequestID())
	app.UseJWTMiddleware(service, jwtMiddle)
	app.UseKeyMiddleware(service, keyMiddle)
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())
	service.Use(keyMiddle)

	// Mount "session" controller
	c4 := NewSessionController(service, jwtSec)
	app.MountSessionController(service, c4)
	// Mount "facebook" controller
	c := NewFacebookController(service, jwtSec, c4)
	app.MountFacebookController(service, c)
	// Mount "google" controller
	c2 := NewGoogleController(service, jwtSec, c4)
	app.MountGoogleController(service, c2)
	// Mount "password-auth" controller

	l := NewLinkedinController(service, jwtSec, c4)
	app.MountLinkedinController(service, l)

	m := NewMicrosoftController(service, jwtSec, c4)
	app.MountMicrosoftController(service, m)

	t := NewTwitterController(service, jwtSec, c4)
	app.MountTwitterController(service, t)

	a := NewAmazonController(service, jwtSec, c4)
	app.MountAmazonController(service, a)

	n := NewNewsletterController(service, jwtSec)
	app.MountNewsletterController(service, n)

	b := NewBookmarkController(service, jwtSec)
	app.MountBookmarkController(service, b)

	c3 := NewPasswordAuthController(service, jwtSec, c4)
	app.MountPasswordAuthController(service, c3)
	// Mount "user" controller
	c5 := NewUserController(service, jwtSec)
	app.MountUserController(service, c5)

	log2.Fatal(service.ListenAndServe("http://localhost:4000"))
}

func initDB() {
	if database.Database == nil {
		database.InitDB()
	}

	models.InitCollections()
}
