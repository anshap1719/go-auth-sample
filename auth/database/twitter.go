package database

import (
	"context"
	"errors"
	"gigglesearch.org/giggle-auth/auth/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gofrs/uuid"
	"github.com/mrjones/oauth"
	"time"
)

var ErrTwitterAccountNotFound = errors.New("No TwitterAccount found in the database")
var ErrTwitterConnectionNotFound = errors.New("No TwitterConnection found in the database")
var ErrTwitterRegisterNotFound = errors.New("No TwitterRegister found in the database")

type TwitterRegister struct {
	TwitterID string `bson:"facebook_id"`

	ID uuid.UUID `bson:"id"`

	TimeCreated time.Time `bson:"time_created"`
}

type TwitterConnection struct {
	MergeToken uuid.UUID `bson:"merge_token"`

	Purpose int `bson:"purpose"`

	State uuid.UUID `bson:"state"`

	TimeCreated time.Time `bson:"time_created"`
}

type TwitterAccount struct {
	ID string `bson:"id"`

	UserID string `bson:"user_id"`
}

func CreateTwitterAccount(ctx context.Context, newTwitterAccount *TwitterAccount) (err error) {
	if err := models.TwitterAccountCollection.Insert(newTwitterAccount); err != nil {
		return err
	}

	return nil
}

func GetTwitterAccount(ctx context.Context, ID string) (*TwitterAccount, error) {
	var fb TwitterAccount

	if err := models.TwitterAccountCollection.Find(bson.M{"id": ID}).One(&fb); err == mgo.ErrNotFound {
		return nil, ErrTwitterAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &fb, nil
}

func DeleteTwitterAccount(ctx context.Context, ID string) error {
	return models.TwitterAccountCollection.Remove(bson.M{"id": ID})
}

func QueryTwitterAccountUser(ctx context.Context, UserID string) (string, error) {
	var fb TwitterAccount

	if err := models.TwitterAccountCollection.Find(bson.M{"user_id": UserID}).One(&fb); err == mgo.ErrNotFound {
		return "", ErrTwitterAccountNotFound
	} else if err != nil {
		return "", err
	}

	return fb.ID, nil
}

func CreateTwitterConnection(ctx context.Context, newTwitterConnection *TwitterConnection) (State uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newTwitterConnection.State = uid

	if err := models.TwitterConnectionCollection.Insert(newTwitterConnection); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetTwitterConnection(ctx context.Context, State uuid.UUID) (*TwitterConnection, error) {
	var fb TwitterConnection

	if err := models.TwitterConnectionCollection.Find(bson.M{"state": State}).One(&fb); err == mgo.ErrNotFound {
		return nil, ErrTwitterConnectionNotFound
	} else if err != nil {
		return nil, err
	}

	return &fb, nil
}

func DeleteTwitterConnection(ctx context.Context, State uuid.UUID) error {
	return models.TwitterConnectionCollection.Remove(bson.M{"state": State})
}

func CreateTwitterRegister(ctx context.Context, newTwitterRegister *TwitterRegister) (ID uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newTwitterRegister.ID = uid

	if err := models.TwitterRegisterCollection.Insert(newTwitterRegister); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetTwitterRegister(ctx context.Context, ID uuid.UUID) (*TwitterRegister, error) {
	var fb TwitterRegister

	if err := models.TwitterRegisterCollection.Find(bson.M{"id": ID}).One(&fb); err == mgo.ErrNotFound {
		return nil, ErrTwitterRegisterNotFound
	} else if err != nil {
		return nil, err
	}

	return &fb, nil
}

func DeleteTwitterRegister(ctx context.Context, ID uuid.UUID) error {
	return models.TwitterRegisterCollection.Remove(bson.M{"id": ID})
}

func GetTwitterToken(key string) (*oauth.RequestToken, error) {
	var token struct {
		Key   string
		Token oauth.RequestToken
	}

	if err := models.TwitterTokenCollection.Find(bson.M{"key": key}).One(&token); err != nil {
		return nil, err
	}

	return &token.Token, nil
}

func CreateTwitterToken(key string, token oauth.RequestToken) error {
	var t = struct {
		Key   string
		Token oauth.RequestToken
	}{
		Key:   key,
		Token: token,
	}

	return models.TwitterTokenCollection.Insert(&t)
}

func DeleteTwitterToken(key string) error {
	return models.TwitterTokenCollection.Remove(bson.M{"key": key})
}
