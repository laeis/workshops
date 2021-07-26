package entities

import (
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
	"workshops/rest-api/internal/config"
)

//JWT token presenter
type JWT struct {
	Token string `json:"access_token"`
}

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	Email string
	jwt.StandardClaims
}

func NewJwtWrapper(c config.AppConfig) JwtWrapper {
	return JwtWrapper{
		SecretKey:       c.SecretKey,
		Issuer:          "AuthService",
		ExpirationHours: c.ExpirationHours,
	}
}

// GenerateToken generates a jwt token
func (j *JwtWrapper) GenerateToken(email string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

//ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return nil, err
	}

	return claims, nil

}

func (j *JwtWrapper) ParseToken(r *http.Request) (clientToken string, err error) {
	clientToken = r.Header.Get("Authorization")
	if clientToken == "" {
		err = errors.New("Empty 'Authorization' header")
		return
	}
	extractedToken := strings.Split(clientToken, "Bearer ")

	if len(extractedToken) == 2 {
		clientToken = strings.TrimSpace(extractedToken[1])
	} else {
		err = errors.New("Wrong token content")
		return
	}
	return
}
