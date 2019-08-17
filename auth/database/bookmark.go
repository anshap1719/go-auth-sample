package database

import (
	"context"
	"errors"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/models"
)

func AddBookmark(userID string, bookmark models.Bookmark) error {
	user, err := GetUser(context.Background(), userID)
	if err != nil {
		return err
	}

	var exists = false

	for _, b := range user.Bookmarks {
		if b.ID == bookmark.ID {
			exists = true
		}
	}

	if exists {
		return errors.New("this resource is already bookmarked")
	} else {
		user.Bookmarks = append(user.Bookmarks, bookmark)
	}

	if err := UpdateUser(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func RemoveFromBookmark(userID, id string) error {
	user, err := GetUser(context.Background(), userID)
	if err != nil {
		return err
	}

	j := 0

	for _, n := range user.Bookmarks {
		if n.ID == id {
			user.Bookmarks[j] = n
			j++
		}
	}

	user.Bookmarks = user.Bookmarks[:j]

	if err := UpdateUser(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func GetAllBookmarks(userID string) ([]models.Bookmark, error) {
	user, err := GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return user.Bookmarks, nil
}

func GetPostBookmarks(userID string) ([]models.Bookmark, error) {
	user, err := GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	j := 0

	for _, n := range user.Bookmarks {
		if n.Type == "post" {
			user.Bookmarks[j] = n
			j++
		}
	}

	user.Bookmarks = user.Bookmarks[:j]

	return user.Bookmarks, nil
}

func GetVideoBookmarks(userID string) ([]models.Bookmark, error) {
	user, err := GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	j := 0

	for _, n := range user.Bookmarks {
		if n.Type == "video" {
			user.Bookmarks[j] = n
			j++
		}
	}

	user.Bookmarks = user.Bookmarks[:j]

	return user.Bookmarks, nil
}

func BookmarkToAppBookmark(b models.Bookmark) *app.Bookmark {
	return &app.Bookmark{
		ID:          b.ID,
		Type:        b.Type,
		Title:       b.Title,
		Description: b.Description,
		Category:    b.Category,
	}
}

func BookmarksToAppBookmarks(marks []models.Bookmark) []*app.Bookmark {
	var appBookmarks []*app.Bookmark

	for _, b := range marks {
		ab := BookmarkToAppBookmark(b)
		appBookmarks = append(appBookmarks, ab)
	}

	return appBookmarks
}

func AppBookmarkToBookmark(b *app.BookmarkParams) models.Bookmark {
	return models.Bookmark{
		ID:          *b.ID,
		Title:       *b.Title,
		Type:        *b.Type,
		Category:    *b.Category,
		Description: *b.Description,
	}
}
