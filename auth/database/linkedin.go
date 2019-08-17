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

var ErrLinkedinAccountNotFound = errors.New("No LinkedinAccount found in the database")
var ErrLinkedinConnectionNotFound = errors.New("No LinkedinConnection found in the database")
var ErrLinkedinRegisterNotFound = errors.New("No LinkedinRegister found in the database")

type LinkedinRegister struct {
	ID uuid.UUID `bson:"id"`

	LinkedinEmail string `bson:"linkedin_email"`

	TimeCreated time.Time `bson:"time_created"`
}

type LinkedinConnection struct {
	MergeToken uuid.UUID `bson:"merge_token"`

	Purpose int `bson:"purpose"`

	State uuid.UUID `bson:"state"`

	TimeCreated time.Time `bson:"time_created"`
}

type LinkedinAccount struct {
	LinkedinEmail string `bson:"linkedin_email"`

	UserID string `bson:"user_id"`
}

func CreateLinkedinAccount(ctx context.Context, newLinkedinAccount *LinkedinAccount) (err error) {
	return models.LinkedinAccountCollection.Insert(newLinkedinAccount)
}

func GetLinkedinAccount(ctx context.Context, LinkedinEmail string) (*LinkedinAccount, error) {
	var la LinkedinAccount

	if err := models.LinkedinAccountCollection.Find(bson.M{"linkedin_email": LinkedinEmail}).One(&la); err == mgo.ErrNotFound {
		return nil, ErrLinkedinAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &la, nil
}

func DeleteLinkedinAccount(ctx context.Context, LinkedinEmail string) error {
	return models.LinkedinAccountCollection.Remove(bson.M{"linkedin_email": LinkedinEmail})
}

func QueryLinkedinAccountUser(ctx context.Context, UserID string) (string, error) {
	var ID []string

	if err := models.LinkedinAccountCollection.Find(bson.M{"user_id": UserID}).Select(bson.M{"linkedin_email": 1}).All(&ID); err == mgo.ErrNotFound {
		return "", ErrLinkedinAccountNotFound
	} else if err != nil {
		return "", err
	}

	if len(ID) == 0 {
		return "", ErrLinkedinAccountNotFound
	}

	return ID[0], nil
}

func CreateLinkedinConnection(ctx context.Context, newLinkedinConnection *LinkedinConnection) (State uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newLinkedinConnection.State = uid

	return uid, models.LinkedinConnectionCollection.Insert(newLinkedinConnection)
}

func GetLinkedinConnection(ctx context.Context, State uuid.UUID) (*LinkedinConnection, error) {
	var lc LinkedinConnection

	if err := models.LinkedinConnectionCollection.Find(bson.M{"state": State}).One(&lc); err == mgo.ErrNotFound {
		return nil, ErrLinkedinConnectionNotFound
	} else if err != nil {
		return nil, err
	}

	return &lc, nil
}

func DeleteLinkedinConnection(ctx context.Context, State uuid.UUID) error {
	return models.LinkedinConnectionCollection.Remove(bson.M{"state": State})
}

func CreateLinkedinRegister(ctx context.Context, newLinkedinRegister *LinkedinRegister) (ID uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newLinkedinRegister.ID = uid

	return uid, models.LinkedinRegisterCollection.Insert(newLinkedinRegister)
}

func GetLinkedinRegister(ctx context.Context, ID uuid.UUID) (*LinkedinRegister, error) {
	var lr LinkedinRegister

	if err := models.LinkedinRegisterCollection.Find(bson.M{"id": ID}).One(&lr); err == mgo.ErrNotFound {
		return nil, ErrLinkedinRegisterNotFound
	} else if err != nil {
		return nil, err
	}

	return &lr, nil
}

func DeleteLinkedinRegister(ctx context.Context, ID uuid.UUID) error {
	return models.LinkedinRegisterCollection.Remove(bson.M{"id": ID})
}
