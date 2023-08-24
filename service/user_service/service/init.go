/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-21 02:01:31
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-22 16:13:40
 * @FilePath: \Road2TikTok\service\user_service\service\init.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */

package user_service

import "github.com/Road-To-Byte/Road2TikTok/pkg/jwt"

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}
