package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("facebook", func() {
	BasePath("/auth/facebook")

	Action("login", func() {
		Description("Gets the URL the front-end should redirect the browser to in order to be authenticated with Facebook, to be logged in")
		Routing(GET("/login"))
		Params(func() {
			Param("token", UUID, "A merge token for merging into an account")
		})

		Response(OK, "text/plain")
		Response(InternalServerError, ErrorMedia)
	})

	Action("register-url", func() {
		Description("Gets the URL the front-end should redirect the browser to in order to be authenticated with Facebook, and then register")
		Routing(GET("/register-start"))

		Response(OK, "text/plain")
		Response(InternalServerError, ErrorMedia)
	})

	Action("register", func() {
		Description("Registers a new account with the system, with Facebook as the login system")
		Routing(POST("/register"))
		Payload(FacebookRegisterParams)

		Response(OK, UserMedia, func() {
			Headers(func() {
				Header("Authorization")
				Header("X-Session")
				Required("Authorization", "X-Session")
			})
		})
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("attach-to-account", func() {
		Description("Attaches a Facebook account to an existing user account, returns the URL the browser should be redirected to")
		Routing(POST("/attach"))
		Security(JWTSec)

		Response(OK, "text/plain")
		Response(InternalServerError, ErrorMedia)
	})

	Action("detach-from-account", func() {
		Description("Detaches a Facebook account from an existing user account.")
		Routing(POST("/detach"))
		Security(JWTSec)

		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("receive", func() {
		Description("The endpoint that Facebook redirects the browser to after the user has authenticated")
		Routing(GET("/receive"))
		Params(func() {
			Param("code", String)
			Param("state", UUID)
			Required("code", "state")
		})

		Response(OK, Any)
		Response(Unauthorized, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})
