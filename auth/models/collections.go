package models

import (
	"gigglesearch.org/giggle-auth/utils/database"
	"github.com/globalsign/mgo"
)

var PasswordLoginCollection *mgo.Collection
var UsersCollection *mgo.Collection
var ResetPasswordCollection *mgo.Collection
var FacebookAccountCollection *mgo.Collection
var FacebookRegisterCollection *mgo.Collection
var FacebookConnectionCollection *mgo.Collection
var TwitterAccountCollection *mgo.Collection
var TwitterRegisterCollection *mgo.Collection
var TwitterConnectionCollection *mgo.Collection
var TwitterTokenCollection *mgo.Collection
var GoogleAccountCollection *mgo.Collection
var GoogleRegisterCollection *mgo.Collection
var GoogleConnectionCollection *mgo.Collection
var AmazonAccountCollection *mgo.Collection
var AmazonRegisterCollection *mgo.Collection
var AmazonConnectionCollection *mgo.Collection
var LinkedinAccountCollection *mgo.Collection
var LinkedinRegisterCollection *mgo.Collection
var LinkedinConnectionCollection *mgo.Collection
var MicrosoftAccountCollection *mgo.Collection
var MicrosoftRegisterCollection *mgo.Collection
var MicrosoftConnectionCollection *mgo.Collection
var NewsletterSubscribersCollection *mgo.Collection
var SessionsCollection *mgo.Collection
var LoginTokenCollection *mgo.Collection
var MergeTokenCollection *mgo.Collection
var EmailVerificationCollection *mgo.Collection

func InitCollections() {
	PasswordLoginCollection = database.GetCollection("password-login")
	UsersCollection = database.GetCollection("users")
	ResetPasswordCollection = database.GetCollection("reset-password")
	FacebookAccountCollection = database.GetCollection("facebook-account")
	FacebookConnectionCollection = database.GetCollection("facebook-connection")
	FacebookRegisterCollection = database.GetCollection("facebook-register")
	TwitterAccountCollection = database.GetCollection("twitter-account")
	TwitterConnectionCollection = database.GetCollection("twitter-connection")
	TwitterRegisterCollection = database.GetCollection("twitter-register")
	TwitterTokenCollection = database.GetCollection("twitter-token")
	GoogleAccountCollection = database.GetCollection("google-account")
	GoogleConnectionCollection = database.GetCollection("google-connection")
	GoogleRegisterCollection = database.GetCollection("google-register")
	AmazonAccountCollection = database.GetCollection("amazon-account")
	AmazonConnectionCollection = database.GetCollection("amazon-connection")
	AmazonRegisterCollection = database.GetCollection("amazon-register")
	LinkedinAccountCollection = database.GetCollection("linkedin-account")
	LinkedinConnectionCollection = database.GetCollection("linkedin-connection")
	LinkedinRegisterCollection = database.GetCollection("linkedin-register")
	MicrosoftAccountCollection = database.GetCollection("microsoft-account")
	MicrosoftConnectionCollection = database.GetCollection("microsoft-connection")
	MicrosoftRegisterCollection = database.GetCollection("microsoft-register")
	NewsletterSubscribersCollection = database.GetCollection("newsletter-subscribers")
	SessionsCollection = database.GetCollection("sessions")
	LoginTokenCollection = database.GetCollection("login-token")
	MergeTokenCollection = database.GetCollection("merge-token")
	EmailVerificationCollection = database.GetCollection("email-verification")
}
