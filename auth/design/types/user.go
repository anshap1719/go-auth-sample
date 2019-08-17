package types

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

const (
	minNameLength = 2
	maxNameLength = 50
)

var UserMedia = MediaType("user", func() {
	Description("A user in the system")
	ContentType("application/json")
	Attributes(func() {
		Attribute("id", String, "Unique unchanging user ID", func() {
			Metadata("struct:field:type", "string")
		})
		Attribute("firstName", String, "Given name for the user", func() {
			Example("Jeff")
		})
		Attribute("lastName", String, "Family name for the user", func() {
			Example("Newmann")
		})
		Attribute("email", String, "Email attached to the account of the user")
		Attribute("phone", String, "Phone Number Of the user")
		Attribute("category", ArrayOf(String), "Category/Categories that a user might select (User interests)")
		Attribute("gender", String)
		Attribute("profileImage", String)
		Attribute("bookmarks", ArrayOf(BookmarkMedia))
		Attribute("changingEmail", String, "When the user attempts to change their email, this is what they will change it to after they verify that it belongs to them")
		Attribute("verifiedEmail", Boolean, "Whether the user has verified their email")
		Attribute("isAdmin", Boolean, "Whether the user is an administrator on the site")
		Attribute("isPluginAuthor", Boolean, "Whether the user is a plugin author on the site")
		Attribute("isEventAuthor", Boolean, "Whether the user is a event author on the site")
		Attribute("getNewsletter", Boolean, "True if the user wants to receive the newsletter")
		Attribute("plugins", ArrayOf(UserPluginMedia), "IDs of all plugins that the user has installed")

		Required("id", "email", "phone", "verifiedEmail", "isAdmin", "isPluginAuthor", "firstName", "lastName", "getNewsletter", "profileImage", "gender")
	})
	View("default", func() {
		Attribute("id")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("category")
		Attribute("isPluginAuthor")
		Attribute("isEventAuthor")
		Attribute("profileImage")
		Attribute("plugins")
	})
	View("owner", func() {
		Attribute("id")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("category")
		Attribute("email")
		Attribute("phone")
		Attribute("changingEmail")
		Attribute("verifiedEmail")
		Attribute("getNewsletter")
		Attribute("profileImage")
		Attribute("isPluginAuthor")
		Attribute("isEventAuthor")
		Attribute("gender")
		Attribute("plugins")
	})
	View("admin", func() {
		Attribute("id")
		Attribute("firstName")
		Attribute("lastName")
		Attribute("category")
		Attribute("email")
		Attribute("phone")
		Attribute("profileImage")
		Attribute("gender")
		Attribute("changingEmail")
		Attribute("verifiedEmail")
		Attribute("isAdmin")
		Attribute("isPluginAuthor")
		Attribute("isEventAuthor")
		Attribute("getNewsletter")
		Attribute("plugins")
	})
})

var MultiUserGetParams = ArrayOf(String)

var UpdateUserParams = Type("user-params", func() {
	Attribute("firstName", String, "Given name for the user", func() {
		Example("Jeff")
		MinLength(minNameLength)
		MaxLength(maxNameLength)
	})
	Attribute("lastName", String, "Family name for the user", func() {
		Example("Newmann")
		MinLength(minNameLength)
		MaxLength(maxNameLength)
	})
	Attribute("email", String, "The primary email of the user", func() {
		Format("email")
	})
	Attribute("phone", String)
	Attribute("gender", String)
	Attribute("profileImage", String)
	Attribute("category", ArrayOf(String), "Category/Categories that a user might select (User interests)")
	Attribute("getNewsletter", Boolean, "True if the user wants to receive the newsletter")
	Attribute("isPluginAuthor", Boolean, "True if user is a plugin author")
	Attribute("isEventAuthor", Boolean, "True if user is a event author")
})

var UserAdminUpdateParams = Type("user-params-admin", func() {
	Attribute("getNewsletter", Boolean)
	Attribute("isAdmin", Boolean)
	Attribute("isPluginAuthor", Boolean)
	Attribute("isEventAuthor", Boolean)
	Attribute("verifiedEmail", Boolean)
})

var AuthMedia = MediaType("auth-status", func() {
	Description("If other Oauths or Auths exists on account.")
	ContentType("application/json")
	Attributes(func() {
		Attribute("google", Boolean, "True if user has google Oauth signin")
		Attribute("facebook", Boolean, "True if user has facebook Oauth signin")
		Attribute("twitter", Boolean, "True if user has twitter Oauth signin")
		Attribute("linkedin", Boolean, "True if user has linkedin Oauth signin")
		Attribute("microsoft", Boolean, "True if user has microsoft Oauth signin")
		Attribute("standard", Boolean, "True if user has password signin")
		Required("google", "facebook", "twitter", "linkedin", "microsoft", "standard")
	})
	View("default", func() {
		Attribute("google")
		Attribute("facebook")
		Attribute("twitter")
		Attribute("linkedin")
		Attribute("microsoft")
		Attribute("standard")
	})
})

var UserPluginMedia = MediaType("user-plugin-media", func() {
	Description("Details about all the plugins that user has installed")
	ContentType("application/json")

	Attributes(func() {
		Attribute("pluginID", String)
		Attribute("permissionsAllowed", ArrayOf(String))
		Attribute("dateAdded", DateTime)

		Required("pluginID", "permissionsAllowed", "dateAdded")
	})

	View("default", func() {
		Attribute("pluginID")
		Attribute("permissionsAllowed")
		Attribute("dateAdded")
	})
})

var UserPlugin = Type("user-plugin", func() {
	Attribute("pluginID", String)
	Attribute("permissionsAllowed", ArrayOf(String))
})
