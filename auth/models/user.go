package models

import (
	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName      string        `json:"first_name" bson:"first_name"`
	LastName       string        `json:"last_name" bson:"last_name"`
	Email          string        `json:"email" bson:"email"`
	Phone          string        `json:"phone" bson:"phone"`
	Category       []string      `json:"category" bson:"category"`
	Bookmarks      []Bookmark    `json:"bookmarks" bson:"bookmarks"`
	ChangingEmail  string        `json:"changing_email" bson:"changing_email"`
	VerifiedEmail  bool          `json:"verified_email" bson:"verified_email"`
	IsAdmin        bool          `json:"is_admin" bson:"is_admin"`
	IsPluginAuthor bool          `json:"is_plugin_author" bson:"is_plugin_author"`
	IsEventAuthor  bool          `json:"is_event_author" bson:"is_event_author"`
	GetNewsletter  bool          `json:"get_newsletter" bson:"get_newsletter"`
	Password       string        `json:"password" bson:"password"`
	ProfileImage   string        `json:"profile_image" bson:"profile_image"`
	Gender         string        `json:"gender" bson:"gender"`
}

type PasswordLogin struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Email string `bson:"email"`

	Password string `bson:"password"`

	Recovery string `bson:"recovery"`

	UserID string `bson:"user_id"`
}
