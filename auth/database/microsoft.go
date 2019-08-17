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

var ErrMicrosoftAccountNotFound = errors.New("No MicrosoftAccount found in the database")
var ErrMicrosoftConnectionNotFound = errors.New("No MicrosoftConnection found in the database")
var ErrMicrosoftRegisterNotFound = errors.New("No MicrosoftRegister found in the database")

type MicrosoftRegister struct {
	ID uuid.UUID `bson:"id"`

	MicrosoftEmail string `bson:"microsoft_email"`

	TimeCreated time.Time `bson:"time_created"`
}

type MicrosoftConnection struct {
	MergeToken uuid.UUID `bson:"merge_token"`

	Purpose int `bson:"purpose"`

	State uuid.UUID `datastore:"-"`

	TimeCreated time.Time `bson:"time_created"`
}

type MicrosoftAccount struct {
	MicrosoftEmail string `bson:"microsoft_email"`

	UserID string `bson:"user_id"`
}

func CreateMicrosoftAccount(ctx context.Context, newMicrosoftAccount *MicrosoftAccount) (err error) {
	return models.MicrosoftAccountCollection.Insert(newMicrosoftAccount)
}

func GetMicrosoftAccount(ctx context.Context, MicrosoftEmail string) (*MicrosoftAccount, error) {
	var ma MicrosoftAccount

	if err := models.MicrosoftAccountCollection.Find(bson.M{"microsoft_email": MicrosoftEmail}).One(&ma); err == mgo.ErrNotFound {
		return nil, ErrMicrosoftAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &ma, nil
}

func DeleteMicrosoftAccount(ctx context.Context, MicrosoftEmail string) error {
	return models.MicrosoftAccountCollection.Remove(bson.M{"microsoft_email": MicrosoftEmail})
}

func QueryMicrosoftAccountUser(ctx context.Context, UserID string) (string, error) {
	var ID []string

	if err := models.MicrosoftAccountCollection.Find(bson.M{"user_id": UserID}).Select(bson.M{"microsoft_email": 1}).All(&ID); err == mgo.ErrNotFound {
		return "", ErrMicrosoftAccountNotFound
	} else if err != nil {
		return "", err
	}

	if len(ID) == 0 {
		return "", ErrMicrosoftAccountNotFound
	}

	return ID[0], nil
}

func CreateMicrosoftConnection(ctx context.Context, newMicrosoftConnection *MicrosoftConnection) (State uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newMicrosoftConnection.State = uid

	return uid, models.MicrosoftConnectionCollection.Insert(newMicrosoftConnection)
}

func GetMicrosoftConnection(ctx context.Context, State uuid.UUID) (*MicrosoftConnection, error) {
	var mc MicrosoftConnection

	if err := models.MicrosoftConnectionCollection.Find(bson.M{"state": State}).One(&mc); err == mgo.ErrNotFound {
		return nil, ErrMicrosoftAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &mc, nil
}

func DeleteMicrosoftConnection(ctx context.Context, State uuid.UUID) error {
	return models.MicrosoftConnectionCollection.Remove(bson.M{"state": State})
}

func CreateMicrosoftRegister(ctx context.Context, newMicrosoftRegister *MicrosoftRegister) (ID uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newMicrosoftRegister.ID = uid

	return uid, models.MicrosoftRegisterCollection.Insert(newMicrosoftRegister)
}

func GetMicrosoftRegister(ctx context.Context, ID uuid.UUID) (*MicrosoftRegister, error) {
	var mr MicrosoftRegister

	if err := models.MicrosoftRegisterCollection.Find(bson.M{"id": ID}).One(&mr); err == mgo.ErrNotFound {
		return nil, ErrMicrosoftAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return &mr, nil
}

func DeleteMicrosoftRegister(ctx context.Context, ID uuid.UUID) error {
	return models.MicrosoftRegisterCollection.Remove(bson.M{"id": ID})
}
