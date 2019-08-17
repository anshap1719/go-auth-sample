package models

import "time"

type NewsletterSubscriber struct {
	Email        string    `json:"email"`
	SubscribedAt time.Time `json:"subscribed_at"`
	IsActive     bool      `json:"is_active"`
}
