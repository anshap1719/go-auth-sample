package database

import (
	"context"
	"errors"
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

var ErrSessionNotFound = errors.New("No Session found in the database")

type Session struct {
	// The browser and browser version connected with this session
	Browser string `bson:"browser"`
	// The latitude and longitude of the last known location of the session
	Coordinates string `bson:"coordinates"`

	ID bson.ObjectId `bson:"_id,omitempty"`
	// The last IP address where this session was used
	IP string `bson:"ip"`

	IsAdmin        bool `bson:"is_admin"`
	IsPluginAuthor bool `bson:"is_plugin_author"`
	IsEventAuthor  bool `bson:"is_event_author"`
	// Whether the session was from a mobile device
	IsMobile bool `bson:"is_mobile"`
	// Time that this session was last used
	LastUsed time.Time `bson:"last_used"`
	// A human-readable string describing the last known location of the session
	Location string `bson:"location"`
	// The OS of the system where this session was used
	Os string `bson:"os"`
	// ID of the user this session is for
	UserID string `bson:"user_id"`
}

func SessionToSession(gen *Session) *app.Session {
	s := &app.Session{
		Browser:     gen.Browser,
		Coordinates: gen.Coordinates,
		ID:          gen.ID.Hex(),
		IP:          gen.IP,
		IsMobile:    gen.IsMobile,
		LastUsed:    gen.LastUsed,
		Location:    gen.Location,
		Os:          gen.Os,
		UserID:      gen.UserID,
	}
	return s
}

func CreateSession(ctx context.Context, newSession *Session) (ID string, err error) {
	id := bson.NewObjectId()
	newSession.ID = id
	if err := models.SessionsCollection.Insert(newSession); err != nil {
		return "", err
	}

	return newSession.ID.Hex(), nil
}

func GetSession(ctx context.Context, ID string) (*Session, error) {
	var t Session
	if err := models.SessionsCollection.Find(bson.M{"_id": bson.ObjectIdHex(ID)}).One(&t); err == mgo.ErrNotFound {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdateSession(ctx context.Context, updatedSession *Session) error {
	if err := models.SessionsCollection.Update(bson.M{"_id": updatedSession.ID}, updatedSession); err == mgo.ErrNotFound {
		return ErrSessionNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func DeleteSession(ctx context.Context, ID string) error {
	if err := models.SessionsCollection.Remove(bson.M{"_id": bson.ObjectIdHex(ID)}); err == mgo.ErrNotFound {
		return ErrSessionNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func DeleteSessionMulti(ctx context.Context, IDs []string) error {
	if len(IDs) == 0 {
		return nil
	}

	var returnErr error = nil

	for _, ID := range IDs {
		err := models.SessionsCollection.Remove(bson.M{"_id": bson.ObjectIdHex(ID)})
		if err != nil {
			returnErr = err
		}
	}

	return returnErr
}

func QuerySessionFromAccount(ctx context.Context, UserID string) ([]*Session, error) {
	var sessions []Session

	if err := models.SessionsCollection.Find(bson.M{"user_id": UserID}).Sort("-last_used").All(&sessions); err == mgo.ErrNotFound {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var data []*Session

	for _, session := range sessions {
		data = append(data, &session)
	}

	return data, nil
}

func QuerySessionIds(ctx context.Context, UserID string) ([]string, error) {
	var sessions []Session

	if err := models.SessionsCollection.Find(bson.M{"user_id": UserID}).Select(bson.M{"_id": 1}).All(&sessions); err == mgo.ErrNotFound {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var IDs []string

	for _, session := range sessions {
		IDs = append(IDs, session.ID.Hex())
	}

	return IDs, nil
}

func QuerySessionOld(ctx context.Context, LastUsed time.Time) ([]string, error) {
	var sessions []bson.ObjectId

	if err := models.SessionsCollection.Find(bson.M{"last_used": bson.M{"$lt": LastUsed}}).Select(bson.M{"_id": 1}).All(&sessions); err == mgo.ErrNotFound {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	var data []string

	for _, session := range sessions {
		data = append(data, session.Hex())
	}

	return data, nil
}
