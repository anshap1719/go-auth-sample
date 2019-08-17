package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("password-auth", func() {
	BasePath("/auth")

	Action("register", func() {
		Description("Register a new user with an email and password")
		Routing(POST("/register"))
		Payload(RegisterParams)
		Response(OK, UserMedia, func() {
			Headers(func() {
				Header("Authorization")
				Header("X-Session")
				Required("Authorization", "X-Session")
			})
		})
		Response(BadRequest, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("login", func() {
		Description("Login a user using an email and password")
		Routing(POST("/login"))
		Params(func() {
			Param("token", UUID, "A merge token for merging into an account")
		})
		Payload(LoginParams)
		Response(OK, UserMedia, func() {
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

	Action("remove", func() {
		Description("Removes using a password as a login method")
		Routing(POST("/remove-password"))
		Security(JWTSec)

		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("change-password", func() {
		Description("Changes the user's current password to a new one, also adds a password to the account if there is none")
		Routing(POST("/change-password"))
		Security(JWTSec)
		Payload(ChangePasswordParams)

		Response(OK, "OK")
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("reset", func() {
		Description("Send an email to user to get a password reset, responds with no content even if the email is not on any user account")
		Routing(POST("/reset-password"))
		Params(func() {
			Param("email", String, "Email of the account to send a password reset", func() {
				Format("email")
			})
			Param("redirect-url", String, "URL to redirect to from the user's email link")
			Required("email", "redirect-url")
		})
		Response(OK, "OK")
		Response(InternalServerError, ErrorMedia)
	})

	Action("confirm-reset", func() {
		Description("Confirms that a reset has been completed and changes the password to the new one passed in")
		Routing(POST("/finalize-reset"))
		Payload(ResetPasswordParams)
		Response(OK, "OK")
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})
