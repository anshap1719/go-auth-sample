package types

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var NewsLetterSubscriber = MediaType("newsletter-subscriber", func() {
	Description("An individual newsletter subscriber")
	ContentType("application/json")
	Attributes(func() {
		Attribute("email", String, func() {
			Format("email")
		})
		Attribute("subscribedAt", DateTime)
		Attribute("isActive", Boolean)

		Required("email", "subscribedAt", "isActive")
	})

	View("default", func() {
		Attribute("email")
		Attribute("subscribedAt")
		Attribute("isActive")
	})
})

var NewsletterParams = Type("newsletter-param", func() {
	Description("Parameters to add/remove a newsletter subscriber")
	Attribute("email", String, func() {
		Format("email")
	})
	Attribute("subscribedAt", DateTime)
	Attribute("isActive", Boolean)
})
