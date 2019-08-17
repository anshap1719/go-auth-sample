package database

import (
	"context"
	"gigglesearch.org/giggle-auth/auth/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

var ErrMergeTokenNotFound = errors.New("No MergeToken found in the database")
var ErrPasswordLoginNotFound = errors.New("No PasswordLogin found in the database")
var ErrResetPasswordNotFound = errors.New("No ResetPassword found in the database")

type ResetPassword struct {
	ID uuid.UUID `bson:"id"`

	TimeExpires time.Time `bson:"time_expires"`

	UserID string `bson:"user_id"`
}

type PasswordLogin struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Email string `bson:"email"`

	Password string `bson:"password"`

	Recovery string `bson:"recovery"`

	UserID string `bson:"user_id"`
}

func GetPasswordLogin(ctx context.Context, Email string) (*PasswordLogin, error) {
	var t PasswordLogin

	if err := models.PasswordLoginCollection.Find(bson.M{"email": Email}).One(&t); err == mgo.ErrNotFound {
		return nil, ErrPasswordLoginNotFound
	} else if err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdatePasswordLogin(ctx context.Context, updatedPasswordLogin *PasswordLogin) error {
	if err := models.PasswordLoginCollection.Update(bson.M{"email": updatedPasswordLogin.Email}, updatedPasswordLogin); err == mgo.ErrNotFound {
		return ErrPasswordLoginNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func DeletePasswordLogin(ctx context.Context, Email string) error {
	if err := models.PasswordLoginCollection.Remove(bson.M{"email": Email}); err == mgo.ErrNotFound {
		return ErrPasswordLoginNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func QueryPasswordLoginFromID(ctx context.Context, UserID string) (string, error) {
	var pl PasswordLogin

	if err := models.UsersCollection.Find(bson.M{"user_id": UserID}).One(&pl); err == mgo.ErrNotFound {
		return "", ErrPasswordLoginNotFound
	} else if err != nil {
		return "", err
	}

	return pl.ID.Hex(), nil
}

func CreateResetPassword(ctx context.Context, newResetPassword *ResetPassword) (err error) {
	return models.ResetPasswordCollection.Insert(newResetPassword)
}

func GetResetPassword(ctx context.Context, UserID string) (*ResetPassword, error) {
	var rp ResetPassword

	if err := models.ResetPasswordCollection.Find(bson.M{"user_id": UserID}).One(&rp); err == mgo.ErrNotFound {
		return nil, ErrResetPasswordNotFound
	} else if err != nil {
		return nil, err
	}

	return &rp, nil
}

func DeleteResetPassword(ctx context.Context, UserID string) error {
	return models.ResetPasswordCollection.Remove(bson.M{"user_id": UserID})
}
