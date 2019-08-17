package types

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

const (
	passwordPattern = `^.*[\w].*$`
	minPassLength   = 6
	maxPassLength   = 100
)

var RegisterParams = Type("register-params", func() {
	Attribute("email", String, "The email that will be attached to the account", func() {
		Format("email")
	})
	Attribute("firstName", String, "The user's given name", func() {
		MinLength(minNameLength)
		MaxLength(maxNameLength)
	})
	Attribute("lastName", String, "The user's family name", func() {
		MaxLength(maxNameLength)
	})
	Attribute("password", String, "The password associated with the new account", func() {
		Pattern(passwordPattern)
		MinLength(minPassLength)
		MaxLength(maxPassLength)
	})
	Attribute("gRecaptchaResponse", String, "The recaptcha response code")
	Attribute("category", ArrayOf(String), "Category/Categories that a user might select (User interests)")
	Required("email", "password", "firstName", "lastName", "gRecaptchaResponse", "category")
})

var LoginParams = Type("login-params", func() {
	Attribute("email", String, "The email address of the account to login to", func() {
		Format("email")
	})
	Attribute("password", String, "The password of the account to login to", func() {
		Pattern(passwordPattern)
		MinLength(minPassLength)
		MaxLength(maxPassLength)
	})
	Attribute("TwoFactor", String, "2 Factor Auth if user has enabled the feature", func() {
		MinLength(6)
		MaxLength(8)
	})
	Required("email", "password")
})

var ChangePasswordParams = Type("change-password-params", func() {
	Attribute("oldPassword", String, "The old password for the current user account", func() {
		Pattern(passwordPattern)
		MinLength(minPassLength)
		MaxLength(maxPassLength)
	})
	Attribute("newPassword", String, "The new password for the current user account", func() {
		Pattern(passwordPattern)
		MinLength(minPassLength)
		MaxLength(maxPassLength)
	})
	Required("oldPassword", "newPassword")
})

var ResetPasswordParams = Type("reset-password-params", func() {
	Attribute("resetCode", String, "The UUID of the password reset, send from the user's email")
	Attribute("userID", String, "The ID of the user to reset the password of")
	Attribute("newPassword", String, "The new password that will be used to login to the account", func() {
		Pattern(passwordPattern)
		MinLength(minPassLength)
		MaxLength(maxPassLength)
	})
	Required("resetCode", "userID", "newPassword")
})
