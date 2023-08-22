/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-15 15:48:03
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-22 14:40:37
 * @FilePath: \Road2TikTok\api_gateway\handler\user.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package handler

import (
	"fmt"

	"github.com/Road-To-Byte/Road2TikTok/api_gateway/rpc"
	"github.com/Road-To-Byte/Road2TikTok/api_gateway/rpc/pb"
	"github.com/gin-gonic/gin"
)

type UserInfoRequestData struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

//	UserInfo	用户信息
//	ctx是上下文 可以接收和回应
func UserInfo(ctx *gin.Context) {
	//  ctx.String(http.StatusOK, "UserInfo")
	//	把发送来的json数据接收到UserInfoData里
	var userInfoRequestData UserInfoRequestData
	if err := ctx.ShouldBindJSON(&userInfoRequestData); err != nil {
		fmt.Println(err.Error())
		// return
	}
	//	TODO:	在这里提前对发送来的json数据做check 如检查UserID是否合法
	//	调用grpc的client 即调用在./rpc/xx.go里的client grpc从server取到服务
	userInfoRequest := &pb.UserInfoRequest{
		// Token:  userInfoRequestData.Token,
		// UserId: userInfoRequestData.UserID,
	}
	rpc.UserInfo(ctx, userInfoRequest)
}
