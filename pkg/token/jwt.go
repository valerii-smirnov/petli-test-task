package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	UserIDClaimName         = "user-id"
	expirationTimeClaimName = "exp"
	issuedAaClimeName       = "iat"
	tokenIDClimeName        = "jti"
)

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(secret string, ttl time.Duration) *JWT {
	return &JWT{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (j JWT) Generate(uid uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		UserIDClaimName:         uid.String(),
		expirationTimeClaimName: jwt.NewNumericDate(time.Now().Add(j.ttl)),
		issuedAaClimeName:       jwt.NewNumericDate(time.Now()),
		tokenIDClimeName:        uuid.New().String(),
	}

	return token.SignedString(j.secret)
}

func (j JWT) Valid(token string) (bool, error) {
	t, err := j.Parse(token)
	if err != nil {
		return false, err
	}

	return t.Valid, nil
}

func (j JWT) Parse(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}
