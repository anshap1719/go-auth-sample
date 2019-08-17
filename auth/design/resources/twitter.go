package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("twitter", func() {
	BasePath("/auth/twitter")

	Action("login", func() {
		Description("Gets the URL the front-end should redirect the browser to in order to be authenticated with Twitter, to be logged in")
		Routing(GET("/login"))
		Params(func() {
			Param("token", UUID, "A merge token for merging into an account")
		})

		Response(OK, "text/plain")
		Response(InternalServerError, ErrorMedia)
	})

	Action("register-url", func() {
		Description("Gets the URL the front-end should redirect the browser to in order to be authenticated with Twitter, and then register")
		Routing(GET("/register-start"))

		Response(OK, "text/plain")
		Response(InternalServerError, ErrorMedia)
	})

	Action("register", func() {
		Description("Registers a new account with the system, with Twitter as the login system")
		Routing(POST("/register"))
		Payload(TwitterRegisterParams)

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
		Description("Attaches a Twitter account to an existing user account, returns the URL the browser should be redirected to")
		Routing(POST("/attach"))
		Security(JWTSec)

		Response(OK, "text/plain")
		Response(InternalServerError, ErrorMedia)
	})

	Action("detach-from-account", func() {
		Description("Detaches a Twitter account from an existing user account.")
		Routing(POST("/detach"))
		Security(JWTSec)

		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("receive", func() {
		Description("The endpoint that Twitter redirects the browser to after the user has authenticated")
		Routing(GET("/receive"))
		Params(func() {
			Param("oauth_token", String)
			Param("oauth_verifier", String)
			Param("state", String)
			Required("oauth_token", "oauth_verifier", "state")
		})

		Response(OK, Any)
		Response(Unauthorized, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})
