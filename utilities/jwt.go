package utilities

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julyskies/gohelpers"

	"file-sharing/constants"
)

type jwtClaims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

func CreateJwt(uid string) (string, error) {
	tokenExpirationSeconds := constants.DEFAULT_JWT_EXPIRATION_SECONDS
	tokenExpirationSecondsString := GetEnv(constants.ENV_NAMES.JwtExpirationSeconds)
	if tokenExpirationSecondsString != "" {
		value, parseError := strconv.Atoi(tokenExpirationSecondsString)
		if parseError != nil {
			tokenExpirationSeconds = value
		}
	}

	tokenSecret := GetEnv(constants.ENV_NAMES.JwtSectet, constants.DEFAULT_JWT_SECRET)

	timestamp := gohelpers.MakeTimestampSeconds()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": timestamp + int64(tokenExpirationSeconds),
			"iat": timestamp,
			"uid": uid,
		},
	)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJwt(token string) (string, error) {
	tokenSecret := GetEnv(constants.ENV_NAMES.JwtSectet, constants.DEFAULT_JWT_SECRET)

	tokenInstance, parseError := jwt.ParseWithClaims(
		token,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if parseError != nil {
		// TODO: handle errors
		return "", parseError
	}
	claims, ok := tokenInstance.Claims.(*jwtClaims)
	if !ok {
		return "", errors.New("invalid token")
	}
	return claims.UID, nil
}
