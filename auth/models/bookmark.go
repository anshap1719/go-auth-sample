package models

type Bookmark struct {
	ID          string `bson:"id"`
	Type        string `bson:"type"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Category    string `bson:"category"`
}
