package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("newsletter", func() {
	BasePath("/newsletter")

	Action("get-subscribers", func() {
		Description("Get All Subscribers")
		Routing(GET("/all"))
		Security(JWTSec)

		Response(OK, CollectionOf(NewsLetterSubscriber))
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("get-subscriber-by-email", func() {
		Description("Get a Subscriber using their email")
		Routing(GET(""))

		Params(func() {
			Param("email", String, func() {
				Format("email")
			})
		})

		Response(OK, NewsLetterSubscriber)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("add-subscriber", func() {
		Description("Add a new newsletter subscriber")
		Routing(POST(""))
		Payload(NewsletterParams)

		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("remove-subscriber", func() {
		Description("Remove a newsletter subscriber")
		Routing(DELETE(""))
		Params(func() {
			Param("email", String, func() {
				Format("email")
			})
		})

		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("update-subscriber", func() {
		Description("Update a newsletter subscriber")
		Routing(PATCH(""))
		Payload(NewsletterParams)

		Response(OK, "OK")
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})
