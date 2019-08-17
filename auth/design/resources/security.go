package resources

import (
	. "github.com/goadesign/goa/design/apidsl"
)

var JWTSec = JWTSecurity("jwt", func() {
	Header("Authorization")
	TokenURL(`/api/auth/login`)
})
