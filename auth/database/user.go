package database

import (
	"context"
	"errors"
	"fmt"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

var ErrUserNotFound = errors.New("No User found in the database")
var ErrEmailVerificationNotFound = errors.New("No EmailVerification found in the database")

type EmailVerification struct {
	Email string `bson:"email"`

	ID string `bson:"id"`

	TimeExpires time.Time `bson:"time_expires"`

	UserID string `bson:"user_id"`
}

func CreateUser(ctx context.Context, newUser *models.User) (ID string, err error) {
	var userID = bson.NewObjectId()
	newUser.ID = userID

	if err := models.UsersCollection.Insert(newUser); err != nil {
		return "", err
	}

	return userID.Hex(), nil
}

func GetUser(ctx context.Context, ID string) (*models.User, error) {
	var t models.User

	fmt.Println(ID)

	if err := models.UsersCollection.FindId(bson.ObjectIdHex(ID)).One(&t); err == mgo.ErrNotFound {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &t, nil
}

func GetUserMulti(ctx context.Context, IDs []string) ([]*models.User, error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	var returnErr error
	var data []*models.User

	for _, ID := range IDs {
		var user models.User

		if err := models.UsersCollection.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&user); err == mgo.ErrNotFound {
			returnErr = ErrUserNotFound
		} else if err != nil {
			returnErr = err
		}

		data = append(data, &user)
	}

	return data, returnErr
}

func GetAllUsers() ([]models.User, error) {
	var user []models.User

	if err := models.UsersCollection.Find(bson.M{}).All(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(ctx context.Context, updatedUser *models.User) error {
	if err := models.UsersCollection.Update(bson.M{"_id": updatedUser.ID}, updatedUser); err == mgo.ErrNotFound {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func CreateEmailVerification(ctx context.Context, newEmailVerification *EmailVerification) error {
	if err := models.EmailVerificationCollection.Insert(newEmailVerification); err != nil {
		return err
	}

	// @TODO: Send Email For Verifying User's Email Using The EmailVerification Data

	return nil
}

func GetEmailVerification(ctx context.Context, ID string) (*EmailVerification, error) {
	var ev EmailVerification

	if err := models.EmailVerificationCollection.Find(bson.M{"id": ID}).One(&ev); err == mgo.ErrNotFound {
		return nil, ErrEmailVerificationNotFound
	} else if err != nil {
		return nil, err
	}

	return &ev, nil
}

func DeleteEmailVerification(ctx context.Context, ID string) error {
	return models.EmailVerificationCollection.Remove(bson.M{"id": ID})
}

func QueryEmailVerificationByUserID(ctx context.Context, UserID string) (string, error) {
	var ev EmailVerification

	if err := models.EmailVerificationCollection.Find(bson.M{"user_id": UserID}).One(&ev); err == mgo.ErrNotFound {
		return "", ErrEmailVerificationNotFound
	} else if err != nil {
		return "", err
	}

	return ev.ID, nil
}

func QueryUserEmail(ctx context.Context, Email string) (*models.User, error) {
	var user models.User

	if err := models.UsersCollection.Find(bson.M{"email": Email}).One(&user); err == mgo.ErrNotFound {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func QueryPasswordLoginFromIDWithBody(ctx context.Context, UserID string) (*PasswordLogin, error) {
	var pl PasswordLogin

	if err := models.PasswordLoginCollection.Find(bson.M{"user_id": UserID}).One(&pl); err == mgo.ErrNotFound {
		return nil, ErrPasswordLoginNotFound
	} else if err != nil {
		return nil, err
	}

	return &pl, nil
}

func UserFromFacebookRegisterParams(gen *app.FacebookRegisterParams) *models.User {
	s := &models.User{
		Email:         gen.Email,
		FirstName:     gen.FirstName,
		LastName:      gen.LastName,
		VerifiedEmail: true,
		Category:      []string{},
	}
	return s
}

func UserFromTwitterRegisterParams(gen *app.TwitterRegisterParams) *models.User {
	s := &models.User{
		Email:         gen.Email,
		FirstName:     gen.FirstName,
		LastName:      gen.LastName,
		VerifiedEmail: true,
		Category:      []string{},
	}
	return s
}

func UserFromGoogleRegisterParams(gen *app.GoogleRegisterParams) *models.User {
	s := &models.User{
		Email:         gen.Email,
		FirstName:     gen.FirstName,
		LastName:      gen.LastName,
		VerifiedEmail: true,
		Category:      []string{},
	}
	return s
}

func UserFromAmazonRegisterParams(gen *app.AmazonRegisterParams) *models.User {
	s := &models.User{
		Email:         gen.Email,
		FirstName:     gen.FirstName,
		LastName:      gen.LastName,
		VerifiedEmail: true,
		Category:      []string{},
	}
	return s
}

func UserFromMicrosoftRegisterParams(gen *app.MicrosoftRegisterParams) *models.User {
	s := &models.User{
		Email:         gen.Email,
		FirstName:     gen.FirstName,
		LastName:      gen.LastName,
		VerifiedEmail: true,
		Category:      []string{},
	}
	return s
}

func UserFromLinkedinRegisterParams(gen *app.LinkedinRegisterParams) *models.User {
	s := &models.User{
		Email:         gen.Email,
		FirstName:     gen.FirstName,
		LastName:      gen.LastName,
		VerifiedEmail: true,
		Category:      []string{},
	}
	return s
}

func UserFromUserParamsMerge(from *app.UserParams, to *models.User) *models.User {
	if from.Email != nil {
		to.Email = *from.Email
	}

	if from.FirstName != nil {
		to.FirstName = *from.FirstName
	}

	if from.GetNewsletter != nil {
		to.GetNewsletter = *from.GetNewsletter
	}

	if from.LastName != nil {
		to.LastName = *from.LastName
	}

	if from.Category != nil {
		to.Category = from.Category
	}

	if from.Phone != nil {
		to.Phone = *from.Phone
	}

	if from.IsPluginAuthor != nil {
		to.IsPluginAuthor = *from.IsPluginAuthor
	}

	if from.IsEventAuthor != nil {
		to.IsEventAuthor = *from.IsEventAuthor
	}

	return to
}

func UserToUserAdmin(gen *models.User) *app.UserAdmin {
	s := &app.UserAdmin{
		ChangingEmail:  &gen.ChangingEmail,
		Email:          gen.Email,
		Phone:          gen.Phone,
		FirstName:      gen.FirstName,
		GetNewsletter:  gen.GetNewsletter,
		ID:             gen.ID.Hex(),
		IsAdmin:        gen.IsAdmin,
		LastName:       gen.LastName,
		VerifiedEmail:  gen.VerifiedEmail,
		Category:       gen.Category,
		IsPluginAuthor: gen.IsPluginAuthor,
		IsEventAuthor:  &gen.IsEventAuthor,
	}
	return s
}

func UserToUser(gen *models.User) *app.User {
	s := &app.User{
		FirstName:      gen.FirstName,
		ID:             gen.ID.Hex(),
		LastName:       gen.LastName,
		Category:       gen.Category,
		IsPluginAuthor: gen.IsPluginAuthor,
		IsEventAuthor:  &gen.IsEventAuthor,
	}
	return s
}

func UserToUserOwner(gen *models.User) *app.UserOwner {
	s := &app.UserOwner{
		ChangingEmail:  &gen.ChangingEmail,
		Email:          gen.Email,
		Phone:          gen.Phone,
		FirstName:      gen.FirstName,
		GetNewsletter:  gen.GetNewsletter,
		ID:             gen.ID.Hex(),
		LastName:       gen.LastName,
		VerifiedEmail:  gen.VerifiedEmail,
		Category:       gen.Category,
		IsPluginAuthor: gen.IsPluginAuthor,
		IsEventAuthor:  &gen.IsEventAuthor,
	}
	return s
}
