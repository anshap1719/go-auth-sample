package types

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var AmazonRegisterMedia = MediaType("amazon-register-media", func() {
	Description("Information used to pre-populate the register page")
	ContentType("application/json")
	Attributes(func() {
		Attribute("email", String, "An email extracted from the Amazon account")
		Attribute("firstName", String, "The given name for the user, pulled from the Amazon account", func() {
			Example("Jeff")
		})
		Attribute("lastName", String, "The family name for the user, pulled from the Amazon account", func() {
			Example("Newmann")
		})
		Attribute("oauthKey", UUID, "A key used to connect the register request with the specific account")
		Required("email", "firstName", "lastName", "oauthKey")
	})
	View("default", func() {
		Attribute("email")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("oauthKey")
	})
})

var AmazonRegisterParams = Type("amazon-register-params", func() {
	Attribute("email", String, "The email that will be connected to the account", func() {
		Format("email")
	})
	Attribute("firstName", String, "The given name for the user", func() {
		MinLength(minNameLength)
		MaxLength(maxNameLength)
	})
	Attribute("lastName", String, "The family name for the user", func() {
		MinLength(minNameLength)
		MaxLength(maxNameLength)
	})
	Attribute("oauthKey", UUID, "The key given when the register was approved")
	Attribute("gRecaptchaResponse", String, "The recaptcha response code")
	Required("email", "firstName", "lastName", "oauthKey", "gRecaptchaResponse")
})
