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

var ErrAmazonConnectionNotFound = errors.New("No AmazonConnection found in the database")
var ErrAmazonAccountNotFound = errors.New("No AmazonAccount found in the database")
var ErrAmazonRegisterNotFound = errors.New("No AmazonRegister found in the database")

type AmazonRegister struct {
	AmazonEmail string `bson:"amazon_email"`

	ID uuid.UUID `bson:"id"`

	TimeCreated time.Time `bson:"time_created"`
}

type AmazonAccount struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	AmazonEmail string `bson:"amazon_email"`

	UserID string `bson:"user_id"`
}

type AmazonConnection struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	MergeToken uuid.UUID `bson:"merge_token"`

	Purpose int `bson:"purpose"`

	State uuid.UUID `bson:"state"`

	TimeCreated time.Time `bson:"time_created"`
}

func CreateAmazonConnection(ctx context.Context, newAmazonConnection *AmazonConnection) (State uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newAmazonConnection.State = uid

	if err := models.AmazonConnectionCollection.Insert(newAmazonConnection); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetAmazonConnection(ctx context.Context, State uuid.UUID) (*AmazonConnection, error) {
	var gc AmazonConnection

	if err := models.AmazonConnectionCollection.Find(bson.M{"state": State}).One(&gc); err == mgo.ErrNotFound {
		return nil, ErrAmazonConnectionNotFound
	} else if err != nil {
		panic(err)
		return nil, err
	}

	return &gc, nil
}

func DeleteAmazonConnection(ctx context.Context, State uuid.UUID) error {
	return models.AmazonConnectionCollection.Remove(bson.M{"state": State})
}

func CreateAmazonAccount(ctx context.Context, newAmazonAccount *AmazonAccount) (err error) {
	return models.AmazonAccountCollection.Insert(newAmazonAccount)
}

func GetAmazonAccount(ctx context.Context, AmazonEmail string) (*AmazonAccount, error) {
	var ga AmazonAccount

	if err := models.AmazonAccountCollection.Find(bson.M{"amazon_email": AmazonEmail}).One(&ga); err == mgo.ErrNotFound {
		return nil, ErrAmazonAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &ga, nil
}

func DeleteAmazonAccount(ctx context.Context, AmazonEmail string) error {
	return models.AmazonAccountCollection.Remove(bson.M{"amazon_email": AmazonEmail})
}

func QueryAmazonAccountUser(ctx context.Context, UserID string) (string, error) {
	var email string

	if err := models.AmazonAccountCollection.Find(bson.M{"user_id": UserID}).Select(bson.M{"amazon_email": 1}).One(&email); err == mgo.ErrNotFound {
		return "", ErrAmazonAccountNotFound
	} else if err != nil {
		return "", err
	}

	return email, nil
}

func CreateAmazonRegister(ctx context.Context, newAmazonRegister *AmazonRegister) (ID uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newAmazonRegister.ID = uid

	if err := models.AmazonRegisterCollection.Insert(newAmazonRegister); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetAmazonRegister(ctx context.Context, ID uuid.UUID) (*AmazonRegister, error) {
	var gr AmazonRegister

	if err := models.AmazonRegisterCollection.Find(bson.M{"id": ID}).One(&gr); err == mgo.ErrNotFound {
		return nil, ErrAmazonRegisterNotFound
	} else if err != nil {
		return nil, err
	}

	return &gr, nil
}

func DeleteAmazonRegister(ctx context.Context, ID uuid.UUID) error {
	return models.AmazonRegisterCollection.Remove(bson.M{"id": ID})
}
