/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-21 01:58:05
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-25 22:27:32
 * @FilePath: \Road2TikTok\service\user_service\main.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Road-To-Byte/Road2TikTok/api_gateway/rpc/pb"
	"github.com/Road-To-Byte/Road2TikTok/pkg/etcd"
	"github.com/Road-To-Byte/Road2TikTok/pkg/viper"
	"github.com/Road-To-Byte/Road2TikTok/pkg/zap"
	user_service "github.com/Road-To-Byte/Road2TikTok/service/user_service/service"
	"google.golang.org/grpc"
)

var (
	config      = viper.Init("user")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey  = config.Viper.GetString("JWT.signingKey")
	logger      = zap.InitLogger()
)

func init() {
	user_service.Init(signingKey)
}

type server struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	fmt.Println(etcd.EtcdAddr)
	serviceTcpAddr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	conn, err := net.ListenTCP("tcp", serviceTcpAddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	grpcS := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcS, &server{})
	//	etcd
	go etcd.RegisterEndPointToEtcd(context.TODO(), serviceAddr, serviceName)
	go func() {
		err := grpcS.Serve(conn)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	}()

}
