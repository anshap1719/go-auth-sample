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

var ErrGoogleConnectionNotFound = errors.New("No GoogleConnection found in the database")
var ErrGoogleAccountNotFound = errors.New("No GoogleAccount found in the database")
var ErrGoogleRegisterNotFound = errors.New("No GoogleRegister found in the database")

type GoogleRegister struct {
	GoogleEmail string `bson:"google_email"`

	ID uuid.UUID `bson:"id"`

	TimeCreated time.Time `bson:"time_created"`
}

type GoogleAccount struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	GoogleEmail string `bson:"google_email"`

	UserID string `bson:"user_id"`
}

type GoogleConnection struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	MergeToken uuid.UUID `bson:"merge_token"`

	Purpose int `bson:"purpose"`

	State uuid.UUID `bson:"state"`

	TimeCreated time.Time `bson:"time_created"`
}

func CreateGoogleConnection(ctx context.Context, newGoogleConnection *GoogleConnection) (State uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newGoogleConnection.State = uid

	if err := models.GoogleConnectionCollection.Insert(newGoogleConnection); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetGoogleConnection(ctx context.Context, State uuid.UUID) (*GoogleConnection, error) {
	var gc GoogleConnection

	if err := models.GoogleConnectionCollection.Find(bson.M{"state": State}).One(&gc); err == mgo.ErrNotFound {
		return nil, ErrGoogleConnectionNotFound
	} else if err != nil {
		panic(err)
		return nil, err
	}

	return &gc, nil
}

func DeleteGoogleConnection(ctx context.Context, State uuid.UUID) error {
	return models.GoogleConnectionCollection.Remove(bson.M{"state": State})
}

func CreateGoogleAccount(ctx context.Context, newGoogleAccount *GoogleAccount) (err error) {
	return models.GoogleAccountCollection.Insert(newGoogleAccount)
}

func GetGoogleAccount(ctx context.Context, GoogleEmail string) (*GoogleAccount, error) {
	var ga GoogleAccount

	if err := models.GoogleAccountCollection.Find(bson.M{"google_email": GoogleEmail}).One(&ga); err == mgo.ErrNotFound {
		return nil, ErrGoogleAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &ga, nil
}

func DeleteGoogleAccount(ctx context.Context, GoogleEmail string) error {
	return models.GoogleAccountCollection.Remove(bson.M{"google_email": GoogleEmail})
}

func QueryGoogleAccountUser(ctx context.Context, UserID string) (string, error) {
	var email string

	if err := models.GoogleAccountCollection.Find(bson.M{"user_id": UserID}).Select(bson.M{"google_email": 1}).One(&email); err == mgo.ErrNotFound {
		return "", ErrGoogleAccountNotFound
	} else if err != nil {
		return "", err
	}

	return email, nil
}

func CreateGoogleRegister(ctx context.Context, newGoogleRegister *GoogleRegister) (ID uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newGoogleRegister.ID = uid

	if err := models.GoogleRegisterCollection.Insert(newGoogleRegister); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetGoogleRegister(ctx context.Context, ID uuid.UUID) (*GoogleRegister, error) {
	var gr GoogleRegister

	if err := models.GoogleRegisterCollection.Find(bson.M{"id": ID}).One(&gr); err == mgo.ErrNotFound {
		return nil, ErrGoogleRegisterNotFound
	} else if err != nil {
		return nil, err
	}

	return &gr, nil
}

func DeleteGoogleRegister(ctx context.Context, ID uuid.UUID) error {
	return models.GoogleRegisterCollection.Remove(bson.M{"id": ID})
}
