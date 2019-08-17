package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("session", func() {
	BasePath("/auth")
	Security(JWTSec)

	Action("refresh", func() {
		Description("Take a user's session token and refresh it, also returns a new authentication token")
		Routing(POST("/session"))
		NoSecurity()
		Headers(func() {
			Header("X-Session", String)
			Required("X-Session")
		})
		Response(OK, func() {
			Headers(func() {
				Header("Authorization")
				Header("X-Session")
				Required("Authorization", "X-Session")
			})
		})
		Response(Unauthorized, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("logout", func() {
		Description("Takes a user's auth token, and logs-out the session associated with it")
		Routing(POST("/logout"))
		Response(OK, "OK")
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("logout-other", func() {
		Description("Logout all sessions for the current user except their current session")
		Routing(POST("/logout/all"))
		Response(OK, "OK")
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("logout-specific", func() {
		Description("Logout of a specific session")
		Routing(POST("/logout/:session-id"))
		Params(func() {
			Param("session-id", String)
			Required("session-id")
		})
		Response(OK, "OK")
		Response(BadRequest, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("get-sessions", func() {
		Description("Gets all of the sessions that are associated with the currently logged in user")
		Routing(GET("/sessions"))
		Response(OK, AllSessionsMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("redeemToken", func() {
		Description("Redeems a login token for credentials")
		NoSecurity()
		Routing(POST("/token"))
		Payload(func() {
			Attribute("token", UUID, "The token to redeem")
			Required("token")
		})

		Response(Created, func() {
			Headers(func() {
				Header("Authorization")
				Header("X-Session")
				Required("Authorization", "X-Session")
			})
		})
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("clean-sessions", func() {
		Description("Deletes all the sessions that have expired")
		Security(JWTSec)
		Routing(GET("/clean/sessions"))
		Response(OK, "OK")
		Response(Forbidden, ErrorMedia)
	})
	Action("clean-login-token", func() {
		Description("Cleans old login tokens from the database")
		Security(JWTSec)
		Routing(GET("/clean/token/login"))
		Response(OK, "OK")
		Response(Forbidden, ErrorMedia)
	})
	Action("clean-merge-token", func() {
		Description("Cleans old account merge tokens from the database")
		Security(JWTSec)
		Routing(GET("/clean/token/merge"))
		Response(OK, "OK")
		Response(Forbidden, ErrorMedia)
	})
})
