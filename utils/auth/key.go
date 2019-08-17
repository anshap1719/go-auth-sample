package auth

import (
	"context"
	"fmt"
	"gigglesearch.org/giggle-auth/utils/database"
	"github.com/globalsign/mgo/bson"
	"github.com/goadesign/goa"
	"net/http"
	"strings"
)

var ErrNoKey = goa.ErrBadRequest("api key must be provided")
var ErrUnauthorized = goa.ErrUnauthorized("invalid api key")

func NewKeyMiddleware() (goa.Middleware, error) {
	validateHandler, err := goa.NewMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		if strings.Contains(r.RequestURI, "/stripe/events") {
			return nil
		}

		key := r.Header.Get("API-Key")

		fmt.Println(key)

		if key == "" {
			return ErrNoKey
		} else {
			if count, err := database.GetCollection("keys").Find(bson.M{"key": key}).Count(); err != nil || count == 0 {
				return ErrUnauthorized
			} else {
				return nil
			}
		}
	})

	return validateHandler, err
}
