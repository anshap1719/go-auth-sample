package database

import (
	"errors"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/models"
	"github.com/globalsign/mgo/bson"
)

var ErrAlreadySubscribed = errors.New("user is already subscribed")

func AddSubscriber(sub models.NewsletterSubscriber) error {
	if count, err := models.NewsletterSubscribersCollection.Find(bson.M{"email": sub.Email}).Count(); err != nil {
		return err
	} else if count > 0 {
		return ErrAlreadySubscribed
	}

	if err := models.NewsletterSubscribersCollection.Insert(&sub); err != nil {
		return err
	}

	return nil
}

func RemoveSubscriber(email string) error {
	return models.NewsletterSubscribersCollection.Remove(bson.M{"email": email})
}

func GetNewsletterSubscribers() ([]models.NewsletterSubscriber, error) {
	var subs []models.NewsletterSubscriber

	if err := models.NewsletterSubscribersCollection.Find(bson.M{}).All(&subs); err != nil {
		return nil, err
	}

	return subs, nil
}

func GetNewsletterSubscriber(email string) (models.NewsletterSubscriber, error) {
	var subs models.NewsletterSubscriber

	if err := models.NewsletterSubscribersCollection.Find(bson.M{"email": email}).One(&subs); err != nil {
		return models.NewsletterSubscriber{}, err
	}

	return subs, nil
}

func UpdateSubscriber(sub models.NewsletterSubscriber) error {
	return models.NewsletterSubscribersCollection.Update(bson.M{"email": sub.Email}, sub)
}

func SubscribersToAppSubscribers(subs []models.NewsletterSubscriber) []*app.NewsletterSubscriber {
	var appSubs []*app.NewsletterSubscriber

	for _, sub := range subs {
		appSubs = append(appSubs, SubscriberToAppSubscriber(sub))
	}

	return appSubs
}

func SubscriberToAppSubscriber(sub models.NewsletterSubscriber) *app.NewsletterSubscriber {
	return &app.NewsletterSubscriber{
		Email:        sub.Email,
		SubscribedAt: sub.SubscribedAt,
		IsActive:     sub.IsActive,
	}
}

func AppSubscriberToSubscriber(sub app.NewsletterParam) models.NewsletterSubscriber {
	return models.NewsletterSubscriber{
		Email:        *sub.Email,
		SubscribedAt: *sub.SubscribedAt,
		IsActive:     *sub.IsActive,
	}
}
