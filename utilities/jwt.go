package utilities

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julyskies/gohelpers"

	"file-sharing/constants"
)

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
