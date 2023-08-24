/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-22 16:03:08
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-22 22:42:48
 * @FilePath: \Road2TikTok\pkg\jwt\jwt.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

//	============ 错误变量 ============
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotValidYet = errors.New("token is not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

//	============ JWT 用于签名的密钥 类型采用byte数组 ============
type JWT struct {
	SigningKey []byte
}

//	============ Claims 即jwt里的声明 自定义jwt载荷中的信息 ============
//	CustomClaims Structured version of Claims Section, as referenced at https://tools.ietf.org/html/rfc7519#section-4.1 See examples for how to use this with your own claim types
type CustomClaims struct {
	Id                 int64
	jwt.StandardClaims //	jwt标准claims
}

//	============ 相关方法 ============

//	构造并返回一个JWT 需要传入一个[]byte类型作为签名密钥
func NewJWT(SigningKey []byte) *JWT {
	return &JWT{
		SigningKey,
	}
}

//	生成一个token 需要传入一个claims
func (this *JWT) GenToken(claims CustomClaims) (string, error) {
	//	使用指定的签名方法创建前面对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//	使用指定的secret前面 获得字符串形式的token
	return token.SignedString(this.SigningKey)
}

//	解析jwt 传入的是字符串形式的token 返回解析后的claims
func (this *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	//	因为是自定义claim结构体 需要用这个方法解析
	//	第二个参数是空的结构体指针 解析后的声明会填进去
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return this.SigningKey, nil
	})
	//	错误处理
	if err != nil {
		//	类型断言 把err接口转成底层的具体类型 ve接收
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 { //	jwt格式错误
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 { //	jwt已过期
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 { //	jwt未生效
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid //	jwt无效
			}

		}
	}
	//	类型断言 成功且jwt有效 就认为解析成功
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	//	解析不出来 就认为是无效
	return nil, ErrTokenInvalid
}
