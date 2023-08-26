/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-21 01:58:05
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-26 17:23:00
 * @FilePath: \Road2TikTok\service\user_service\main.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Road-To-Byte/Road2TikTok/pkg/etcd"
	"github.com/Road-To-Byte/Road2TikTok/pkg/viper"
	"github.com/Road-To-Byte/Road2TikTok/service/user_service/pb"
	user_service "github.com/Road-To-Byte/Road2TikTok/service/user_service/service"
	"google.golang.org/grpc"
)

var (
	config      = viper.Init("user")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	signingKey  = config.Viper.GetString("JWT.signingKey")
)

func init() {
	user_service.Init(signingKey)
}

type server struct {
	pb.UnimplementedUserInfoServiceServer
}

func main() {
	log.Println("user_service start")
	serviceTcpAddr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	// serviceTcpAddr := "127.0.0.1:6661"
	// var err error
	log.Println("1")
	if err != nil {
		log.Fatalln("service地址解析失败：", err.Error())
	}
	conn, err := net.ListenTCP("tcp", serviceTcpAddr)
	// conn, err := net.Listen("tcp", serviceTcpAddr)
	log.Println("2")
	if err != nil {
		log.Fatalln("listen失败：", err.Error())
	}
	grpcS := grpc.NewServer()
	log.Println("3")
	pb.RegisterUserInfoServiceServer(grpcS, &server{})
	log.Println("4")
	//	etcd
	go etcd.RegisterEndPointToEtcd(context.TODO(), serviceAddr, serviceName)
	log.Println("5")
	log.Println("Waiting for client requests...")
	serve_err := grpcS.Serve(conn)
	if serve_err != nil {
		log.Fatalln(err.Error())
	}
	// go func() {
	// 	log.Println("Waiting for client requests...")
	// 	err := grpcS.Serve(conn)
	// 	if err != nil {
	// 		log.Fatalln(err.Error())
	// 		return
	// 	}
	// }()
	log.Println("6")
}
