package resources

import (
	. "gigglesearch.org/giggle-auth/auth/design/types"
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("bookmark", func() {
	BasePath("/bookmark")
	DefaultMedia(BookmarkMedia)

	Action("addVideo", func() {
		Description("add a video to bookmarks")
		Routing(PUT("/video"))
		Security(JWTSec)

		Payload(BookmarkParams)

		Response(OK)
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("addPost", func() {
		Description("add a post to bookmarks")
		Routing(PUT("/post"))
		Security(JWTSec)

		Payload(BookmarkParams)

		Response(OK)
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("removeFromBookmark", func() {
		Description("remove a post or video from bookmarks")
		Routing(DELETE(""))
		Security(JWTSec)

		Params(func() {
			Param("id", String, "ID of the resource to be deleted")
		})

		Response(OK)
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getBookmarks", func() {
		Description("get all bookmarks")
		Routing(GET(""))
		Security(JWTSec)

		Response(OK, ArrayOf(BookmarkMedia))
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getPostBookmarks", func() {
		Description("get all post bookmarks")
		Routing(GET("/post"))
		Security(JWTSec)

		Response(OK, ArrayOf(BookmarkMedia))
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getVideoBookmarks", func() {
		Description("get all video bookmarks")
		Routing(GET("/video"))
		Security(JWTSec)

		Response(OK, ArrayOf(BookmarkMedia))
		Response(Unauthorized, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})
