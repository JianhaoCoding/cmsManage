package helpers

import (
	conf "cms/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(conf.InitConf().JWT.Secret)

type Claims struct {
	UserId uint `json:"adminer_id"`
	jwt.StandardClaims
}

func GenerateToken(userId uint) (string, error) {
	claims := Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			// token-有效期2小时
			ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
