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

var ErrLoginTokenNotFound = errors.New("No Login Token found in the database")

type LoginToken struct {
	TimeExpire time.Time `bson:"time_expire"`

	Token uuid.UUID `bson:"token"`

	UserID string `bson:"user_id"`
}

type MergeToken struct {
	TimeExpire time.Time `bson:"time_expire"`

	Token uuid.UUID `bson:"token"`

	UserID string `bson:"user_id"`
}

func CreateLoginToken(ctx context.Context, newLoginToken *LoginToken) (Token uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newLoginToken.Token = uid

	if err := models.LoginTokenCollection.Insert(newLoginToken); err != nil {
		return uuid.Nil, err
	}

	return newLoginToken.Token, nil
}

func GetLoginToken(ctx context.Context, Token uuid.UUID) (*LoginToken, error) {
	var lt LoginToken

	if err := models.LoginTokenCollection.Find(bson.M{"token": Token}).One(&lt); err == mgo.ErrNotFound {
		return nil, ErrLoginTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return &lt, nil
}

func DeleteLoginToken(ctx context.Context, Token uuid.UUID) error {
	return models.LoginTokenCollection.Remove(bson.M{"token": Token})
}

func DeleteLoginTokenMulti(ctx context.Context, Tokens []uuid.UUID) error {
	if len(Tokens) == 0 {
		return nil
	}

	var returnErr error

	for _, Token := range Tokens {
		if err := models.LoginTokenCollection.Remove(bson.M{"token": Token}); err != nil {
			returnErr = err
		}
	}

	return returnErr
}

func QueryLoginTokenOld(ctx context.Context, TimeExpire time.Time) ([]uuid.UUID, error) {
	var lts []LoginToken

	if err := models.LoginTokenCollection.Find(bson.M{"time_expire": bson.M{"$lt": TimeExpire}}).All(&lts); err == mgo.ErrNotFound {
		return nil, ErrLoginTokenNotFound
	} else if err != nil {
		return nil, err
	}

	var data []uuid.UUID

	for _, token := range lts {
		data = append(data, token.Token)
	}

	return data, nil
}

func CreateMergeToken(ctx context.Context, newMergeToken *MergeToken) (Token uuid.UUID, err error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	newMergeToken.Token = uid

	if err := models.MergeTokenCollection.Insert(newMergeToken); err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func GetMergeToken(ctx context.Context, Token uuid.UUID) (*MergeToken, error) {
	var t MergeToken
	if err := models.MergeTokenCollection.Find(bson.M{"token": Token}).One(&t); err == mgo.ErrNotFound {
		return nil, ErrMergeTokenNotFound
	} else if err != nil {
		return nil, err
	}
	return &t, nil
}

func DeleteMergeToken(ctx context.Context, Token uuid.UUID) error {
	return models.MergeTokenCollection.Remove(bson.M{"token": Token})
}

func DeleteMergeTokenMulti(ctx context.Context, Tokens []uuid.UUID) error {
	if len(Tokens) == 0 {
		return nil
	}

	var returnErr error

	for _, token := range Tokens {
		if err := models.MergeTokenCollection.Remove(bson.M{"token": token}); err != nil {
			returnErr = err
		}
	}

	return returnErr
}

func QueryMergeTokenOld(ctx context.Context, TimeExpire time.Time) ([]uuid.UUID, error) {
	var tokens []uuid.UUID

	if err := models.MergeTokenCollection.Find(bson.M{"time_expire": bson.M{"$lt": TimeExpire}}).Select(bson.M{"token": 1}).All(&tokens); err == mgo.ErrNotFound {
		return nil, ErrMergeTokenNotFound
	} else if err != nil {
		return nil, err
	}

	return tokens, nil
}
