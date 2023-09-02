package favorite_service

import "github.com/Road-To-Byte/Road2TikTok/pkg/jwt"

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}
