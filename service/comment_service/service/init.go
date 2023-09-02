/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-09-02 15:46:02
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-09-02 15:46:23
 * @FilePath: \Road2TikTok\service\comment_service\service\init.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package comment_service

import "github.com/Road-To-Byte/Road2TikTok/pkg/jwt"

var (
	Jwt *jwt.JWT
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
}
