package design

import (
	"gigglesearch.org/giggle-auth/utils/secrets"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("user", func() {
	Title("Users and authentication")
	Description("A service that manages users and authentication to various services")
	Scheme(secrets.Scheme)
	Host(secrets.Hostname)
	Origin("*", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS") // Allow all origins to retrieve the Swagger JSON (CORS)
	})
	BasePath("/api/v1/user")
	Security(APIKeySecurity("key", func() {
		Description("API Key for users API")
		Header("API-Key")
	}))
})
