package ctxdata

import (
	"github.com/golang-jwt/jwt/v4"
)

const Identify = "imooc.com"

// 生成jwt token
func GetJwtToken(secret string, iat, seconds int64, uid string) (string, error) {
	claims := jwt.MapClaims{}
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secret))
}
