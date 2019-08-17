package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gigglesearch.org/giggle-auth/utils/secrets"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/satori/go.uuid"
)

func NewJWTMiddleware(security *goa.JWTSecurity) (goa.Middleware, error) {
	key, err := loadJWTPublicKeys()
	if err != nil {
		return nil, err
	}
	validateHandler, err := goa.NewMiddleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		token := jwt.ContextJWT(ctx)
		claims, ok := token.Claims.(jwtgo.MapClaims)
		if !ok {
			return jwt.ErrJWTError("Not a valid token")
		}
		if !claims.VerifyIssuer("Giggle", true) {
			return jwt.ErrJWTError("Not a valid token")
		}
		if sub, ok := claims["sub"]; !ok {
			numStr, ok := sub.(string)
			if !ok {
				return jwt.ErrJWTError("Not a valid token")
			}
			if num, err := strconv.ParseInt(numStr, 10, 64); err != nil || num <= 0 {
				return jwt.ErrJWTError("Not a valid token")
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return jwt.New(jwt.NewSimpleResolver([]jwt.Key{key}), validateHandler, security), nil
}

type JWTSecurity struct {
	privateKey *rsa.PrivateKey
	publicKey  jwt.Key
}

func NewJWTSecurity() (JWTSecurity, error) {
	privKey, err := jwtgo.ParseRSAPrivateKeyFromPEM([]byte(secrets.JWTPrivateKey))
	if err != nil {
		return JWTSecurity{}, nil
	}
	pubKey, err := loadJWTPublicKeys()
	if err != nil {
		return JWTSecurity{}, nil
	}

	return JWTSecurity{
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

func (j *JWTSecurity) IsAdmin(req *http.Request) bool {
	tokenString := req.Header.Get("Authorization")
	token, err := jwtgo.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return false
	}

	if !claims.VerifyIssuer("Giggle", true) {
		return false
	}

	fmt.Println("Claims: ", claims)

	subject := claims["adm"]
	subStr, ok := subject.(string)
	if !ok {
		return false
	}
	isAdmin, _ := strconv.ParseBool(subStr)
	return isAdmin
}

func (j *JWTSecurity) IsGhost(req *http.Request) bool {
	tokenString := req.Header.Get("Authorization")
	token, err := jwtgo.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return false
	}

	if !claims.VerifyIssuer("Giggle", true) {
		return false
	}

	fmt.Println("Claims: ", claims)

	subject := claims["ghs"]
	subStr, ok := subject.(string)
	if !ok {
		return false
	}
	isGhost, _ := strconv.ParseBool(subStr)
	return isGhost
}

func (j *JWTSecurity) IsPluginAuthor(req *http.Request) bool {
	tokenString := req.Header.Get("Authorization")
	token, err := jwtgo.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return false
	}

	if !claims.VerifyIssuer("Giggle", true) {
		return false
	}

	fmt.Println("Claims: ", claims)

	subject := claims["pla"]
	subStr, ok := subject.(string)
	if !ok {
		return false
	}
	isPluginAuthor, _ := strconv.ParseBool(subStr)
	return isPluginAuthor
}

func (j *JWTSecurity) IsEventAuthor(req *http.Request) bool {
	tokenString := req.Header.Get("Authorization")
	token, err := jwtgo.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return false
	}

	if !claims.VerifyIssuer("Giggle", true) {
		return false
	}

	fmt.Println("Claims: ", claims)

	subject := claims["iea"]
	subStr, ok := subject.(string)
	if !ok {
		return false
	}
	isPluginAuthor, _ := strconv.ParseBool(subStr)
	return isPluginAuthor
}

func (j *JWTSecurity) GetUserID(req *http.Request) string {
	tokenString := req.Header.Get("Authorization")
	token, err := jwtgo.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return ""
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return ""
	}

	if !claims.VerifyIssuer("Giggle", true) {
		return ""
	}

	return claims["sub"].(string)
}

func (j *JWTSecurity) GetSessionCode(req *http.Request) string {
	tokenString := req.Header.Get("X-Session")
	token, err := jwtgo.Parse(tokenString, func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return ""
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return ""
	}

	if !claims.VerifyIssuer("Giggle", true) {
		return ""
	}

	return claims["prn"].(string)
}

func (j *JWTSecurity) GetSessionFromAuth(req *http.Request) string {
	tokenString := req.Header.Get("Authorization")
	token, err := jwtgo.Parse(strings.TrimPrefix(tokenString, "Bearer "), func(t *jwtgo.Token) (interface{}, error) {
		return j.publicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return ""
	}
	claims := token.Claims.(jwtgo.MapClaims)

	if !claims.VerifyIssuer("Giggle", true) {
		return ""
	}

	return claims["ses"].(string)
}

func (j *JWTSecurity) SignSessionToken(expTime time.Duration, sessionID string) (string, error) {
	sesExpTime := time.Now().Add(expTime).Unix()
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS512, jwtgo.MapClaims{
		"iss": "Giggle",
		"exp": sesExpTime,
		"prn": sessionID,
		"iat": time.Now().Unix(),
		"nbf": 2,
	})
	return token.SignedString(j.privateKey)
}

func (j *JWTSecurity) SignAuthToken(expTime time.Duration, sessionID, userID string, isAdmin, isPluginAuthor, isEventAuthor bool) (string, error) {
	authExpTime := time.Now().Add(expTime).Unix()
	tokenID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	fmt.Println("Admin: ", isAdmin)

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS512, jwtgo.MapClaims{
		"iss": "Giggle",
		"exp": authExpTime,
		"jti": tokenID.String(),
		"iat": time.Now().Unix(),
		"nbf": 2,
		"sub": userID,
		"ses": sessionID,
		"adm": strconv.FormatBool(isAdmin),
		"pla": strconv.FormatBool(isPluginAuthor),
		"iea": strconv.FormatBool(isEventAuthor),
		"ghs": strconv.FormatBool(false),
	})
	return token.SignedString(j.privateKey)
}

func loadJWTPublicKeys() (jwt.Key, error) {
	key, err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(secrets.JWTPublicKey))
	if err != nil {
		return nil, err
	}
	return key, nil
}
