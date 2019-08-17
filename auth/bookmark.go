package auth

import (
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/utils/auth"
	"github.com/goadesign/goa"
)

// BookmarkController implements the bookmark resource.
type BookmarkController struct {
	*goa.Controller
	auth.JWTSecurity
}

// NewBookmarkController creates a bookmark controller.
func NewBookmarkController(service *goa.Service, jwtSec auth.JWTSecurity) *BookmarkController {
	return &BookmarkController{
		Controller:  service.NewController("BookmarkController"),
		JWTSecurity: jwtSec,
	}
}

// AddPost runs the addPost action.
func (c *BookmarkController) AddPost(ctx *app.AddPostBookmarkContext) error {
	// BookmarkController_AddPost: start_implement

	bookmark := database.AppBookmarkToBookmark(ctx.Payload)
	bookmark.Type = "post"

	if err := database.AddBookmark(c.GetUserID(ctx.Request), bookmark); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(nil)

	// BookmarkController_AddPost: end_implement
}

// AddVideo runs the addVideo action.
func (c *BookmarkController) AddVideo(ctx *app.AddVideoBookmarkContext) error {
	// BookmarkController_AddVideo: start_implement

	bookmark := database.AppBookmarkToBookmark(ctx.Payload)
	bookmark.Type = "video"

	if err := database.AddBookmark(c.GetUserID(ctx.Request), bookmark); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(nil)

	// BookmarkController_AddVideo: end_implement
}

// GetBookmarks runs the getBookmarks action.
func (c *BookmarkController) GetBookmarks(ctx *app.GetBookmarksBookmarkContext) error {
	// BookmarkController_GetBookmarks: start_implement

	b, err := database.GetAllBookmarks(c.GetUserID(ctx.Request))
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	bookmarks := database.BookmarksToAppBookmarks(b)

	return ctx.OK(bookmarks)

	// BookmarkController_GetBookmarks: end_implement
}

// GetPostBookmarks runs the getPostBookmarks action.
func (c *BookmarkController) GetPostBookmarks(ctx *app.GetPostBookmarksBookmarkContext) error {
	// BookmarkController_GetPostBookmarks: start_implement

	b, err := database.GetPostBookmarks(c.GetUserID(ctx.Request))
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	bookmarks := database.BookmarksToAppBookmarks(b)

	return ctx.OK(bookmarks)

	// BookmarkController_GetPostBookmarks: end_implement
}

// GetVideoBookmarks runs the getVideoBookmarks action.
func (c *BookmarkController) GetVideoBookmarks(ctx *app.GetVideoBookmarksBookmarkContext) error {
	// BookmarkController_GetVideoBookmarks: start_implement

	b, err := database.GetVideoBookmarks(c.GetUserID(ctx.Request))
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	bookmarks := database.BookmarksToAppBookmarks(b)

	return ctx.OK(bookmarks)

	// BookmarkController_GetVideoBookmarks: end_implement
}

// RemoveFromBookmark runs the removeFromBookmark action.
func (c *BookmarkController) RemoveFromBookmark(ctx *app.RemoveFromBookmarkBookmarkContext) error {
	// BookmarkController_RemoveFromBookmark: start_implement

	if err := database.RemoveFromBookmark(c.GetUserID(ctx.Request), *ctx.ID); err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	return ctx.OK(nil)

	// BookmarkController_RemoveFromBookmark: end_implement
}
