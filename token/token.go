package token

import (
	"fmt"
	"github.com/Assyl00/goProject/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	jWTPrivateToken = "SecrteTokenSecrteToken"
	ip              = "192.168.0.107"
)

func GenrateToken(claims *models.JwtClaims, expirationTime time.Time) (string, error) {

	claims.ExpiresAt = expirationTime.Unix()
	claims.IssuedAt = time.Now().UTC().Unix()
	claims.Issuer = ip

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jWTPrivateToken))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString, origin string) (bool, *models.JwtClaims) {
	claims := &models.JwtClaims{}
	token, _ := getTokenFromString(tokenString, claims)
	if token.Valid {
		if e := claims.Valid(); e == nil {
			return true, claims
		}
	}
	return false, claims
}

func GetClaims(tokenString string) models.JwtClaims {
	claims := &models.JwtClaims{}

	_, err := getTokenFromString(tokenString, claims)
	if err == nil {
		return *claims
	}
	return *claims
}
func getTokenFromString(tokenString string, claims *models.JwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jWTPrivateToken), nil
	})
}
