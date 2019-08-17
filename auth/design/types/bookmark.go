package types

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var BookmarkMedia = MediaType("bookmark", func() {
	Description("an individual bookmark")
	ContentType("application/json")

	Attributes(func() {
		Attribute("id", String, "ID of the post or video as per database")
		Attribute("type", String, "whether this bookmark is a video or a post")
		Attribute("title", String, "title of post or video")
		Attribute("category", String, "category of the post or video")
		Attribute("description", String, "description of the post or video")

		Required("id", "type", "title", "category", "description")
	})

	View("default", func() {
		Attribute("id")
		Attribute("type")
		Attribute("title")
		Attribute("category")
		Attribute("description")
	})
})

var BookmarkParams = Type("bookmark-params", func() {
	Attribute("id", String)
	Attribute("type", String)
	Attribute("title", String)
	Attribute("category", String)
	Attribute("description", String)
})
