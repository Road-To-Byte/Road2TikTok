/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-21 01:24:15
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-22 15:42:10
 * @FilePath: \Road2TikTok\api_gateway\rpc\user.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package rpc

import (
	"log"
	"sync"

	"github.com/Road-To-Byte/Road2TikTok/api_gateway/rpc/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	initOnce   sync.Once
	initErr    error
	userClient pb.UserServiceClient
	serverAddr = "127.0.0.1"
	serverIp   = "9999"
)

//	初始化客户端
func InitUserClient() {
	conn, err := grpc.Dial(serverAddr+":"+serverIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed connect rpc server: %v", err)
	}
	userClient = pb.NewUserServiceClient(conn)
}

//	用户信息
func UserInfo(ctx *gin.Context, req *pb.UserInfoRequest) (*pb.UserInfoResponse, error) {
	return userClient.UserInfo(ctx, req)
}

//	用户登录
func Login(ctx *gin.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	return userClient.Login(ctx, req)
}

//	用户注册
func Register(ctx *gin.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	return userClient.Register(ctx, req)
}
