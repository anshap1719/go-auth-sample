package database

import (
	"context"
	"errors"
	"gigglesearch.org/giggle-auth/auth/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gofrs/uuid"
	"time"
)

var ErrFacebookAccountNotFound = errors.New("No FacebookAccount found in the database")
var ErrFacebookConnectionNotFound = errors.New("No FacebookConnection found in the database")
var ErrFacebookRegisterNotFound = errors.New("No FacebookRegister found in the database")

type FacebookRegister struct {
	FacebookID string `bson:"facebook_id"`

	ID uuid.UUID `bson:"id"`

	TimeCreated time.Time `bson:"time_created"`
}

type FacebookConnection struct {
	MergeToken uuid.UUID `bson:"merge_token"`

	Purpose int `bson:"purpose"`

	State uuid.UUID `bson:"state"`

	TimeCreated time.Time `bson:"time_created"`
}

type FacebookAccount struct {
	ID string `bson:"id"`

	UserID string `bson:"user_id"`
}

func CreateFacebookAccount(ctx context.Context, newFacebookAccount *FacebookAccount) (err error) {
	if err := models.FacebookAccountCollection.Insert(newFacebookAccount); err != nil {
		return err
	}

	return nil
}

func GetFacebookAccount(ctx context.Context, ID string) (*FacebookAccount, error) {
	var fb FacebookAccount

	if err := models.FacebookAccountCollection.Find(bson.M{"id": ID}).One(&fb); err == mgo.ErrNotFound {
		return nil, ErrFacebookAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &fb, nil
}

func DeleteFacebookAccount(ctx context.Context, ID string) error {
	return models.FacebookAccountCollection.Remove(bson.M{"id": ID})
}

func QueryFacebookAccountUser(ctx context.Context, UserID string) (string, error) {
	var fb FacebookAccount

	if err := models.FacebookAccountCollection.Find(bson.M{"user_id": UserID}).One(&fb); err == mgo.ErrNotFound {
		return "", ErrFacebookAccountNotFound
	} else if err != nil {
		return "", err
	}

	return fb.ID, nil
}

func CreateFacebookConnection(ctx context.Context, newFacebookConnection *FacebookConnection) (State uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newFacebookConnection.State = uid

	if err := models.FacebookConnectionCollection.Insert(newFacebookConnection); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetFacebookConnection(ctx context.Context, State uuid.UUID) (*FacebookConnection, error) {
	var fb FacebookConnection

	if err := models.FacebookConnectionCollection.Find(bson.M{"state": State}).One(&fb); err == mgo.ErrNotFound {
		return nil, ErrFacebookConnectionNotFound
	} else if err != nil {
		return nil, err
	}

	return &fb, nil
}

func DeleteFacebookConnection(ctx context.Context, State uuid.UUID) error {
	return models.FacebookConnectionCollection.Remove(bson.M{"state": State})
}

func CreateFacebookRegister(ctx context.Context, newFacebookRegister *FacebookRegister) (ID uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newFacebookRegister.ID = uid

	if err := models.FacebookRegisterCollection.Insert(newFacebookRegister); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetFacebookRegister(ctx context.Context, ID uuid.UUID) (*FacebookRegister, error) {
	var fb FacebookRegister

	if err := models.FacebookRegisterCollection.Find(bson.M{"id": ID}).One(&fb); err == mgo.ErrNotFound {
		return nil, ErrFacebookRegisterNotFound
	} else if err != nil {
		return nil, err
	}

	return &fb, nil
}

func DeleteFacebookRegister(ctx context.Context, ID uuid.UUID) error {
	return models.FacebookRegisterCollection.Remove(bson.M{"id": ID})
}
