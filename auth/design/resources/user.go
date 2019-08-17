package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("user", func() {
	BasePath("/user")
	DefaultMedia(UserMedia)

	Action("add-plugin", func() {
		Description("Add a new plugin to user's account")
		Routing(POST("/plugins"))

		Payload(UserPlugin)

		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("retrieve", func() {
		Description("Get user by ID")
		Routing(GET(""))
		Params(func() {
			Param("user-id", String, "The ID of the requested user. If this is not provided, get currently logged in user")
		})
		Response(OK)
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("get-many", func() {
		Description("Get many users by their ID")
		Routing(GET("/multi"))
		Params(func() {
			Param("id", MultiUserGetParams)
		})
		Response(OK, CollectionOf(UserMedia, func() {
			ContentType("application/json")
		}))
		Response(BadRequest, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getAuths", func() {
		Description("Returns whether Oauth is attached or not")
		Routing(GET("/authstat"))
		Params(func() {
			Param("userID", String, "The ID of the requested user. If this is not provide, get currently logged in user")
		})
		Response(OK, AuthMedia)
		Response(BadRequest, ErrorMedia)
		Response(Unauthorized, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("get-by-email", func() {
		Description("Get a user by their email. Only callable by admins")
		Routing(GET("/email"))
		Params(func() {
			Param("email", String, "The email of the requested user.", func() {
				Format("email")
			})
			Required("email")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("update", func() {
		Description("Update a user")
		Routing(PATCH(""))
		Payload(UpdateUserParams)
		Security(JWTSec)
		Response(OK, "OK")
		Response(BadRequest, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("deactivate", func() {
		Description("Disable a user's account")
		Routing(DELETE(""))
		Params(func() {
			Param("id", String, "id of the user to be deactivated when admin is deactivating a user")
			Param("admin", Boolean, "whether admin is requesting this deactivation")
		})
		Security(JWTSec)
		Response(OK, "OK")
		Response(Forbidden, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("validate-email", func() {
		Description("Validates an email address, designed to be called by users directly in their browser")
		Routing(GET("//verifyemail/:validateID"))
		Params(func() {
			Param("validateID", String, "The ID of the validation to be confirmed")
			Required("validateID")
		})
		Response(SeeOther)
		Response(NotFound, func() {
			Media("text/html")
			Status(404)
		})
		Response(InternalServerError, ErrorMedia)
	})

	Action("resend-verify-email", func() {
		Description("Resends a verify email for the current user, also invalidates the link on the previously send email verification")
		Routing(POST("/resend-verify"))
		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("get-all-users", func() {
		Security(JWTSec)
		Description("Get all users")
		Routing(GET("/all"))

		Response(OK, Any)
		Response(BadRequest, ErrorMedia)
		Response(Forbidden, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("update-admin", func() {
		Security(JWTSec)
		Description("Update a user from admin dashboard")
		Routing(PATCH("/update-user"))

		Params(func() {
			Param("uid", String, "user id to be modified")
		})

		Payload(UserAdminUpdateParams)

		Response(OK, "OK")
		Response(Forbidden, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("update-plugin-permissions", func() {
		Description("Update plugin permissions")
		Routing(PUT("/plugins"))

		Payload(UserPlugin)

		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})
