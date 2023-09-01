/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-21 02:01:49
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-09-01 15:41:45
 * @FilePath: \Road2TikTok\service\user_service\service\handler.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package user_service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Road-To-Byte/Road2TikTok/dal/db"
	"github.com/Road-To-Byte/Road2TikTok/pkg/jwt"
	"github.com/Road-To-Byte/Road2TikTok/service/user_service/pb"
)

type UserServiceImpl struct{}

// UserInfo
func (this *UserServiceImpl) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (resp *pb.UserInfoResponse, err error) {
	log.Println("UserInfo get")
	usrID := req.UserId
	//	dal取出user
	usr, err := db.GetUserByUserID(ctx, usrID)
	//	check err
	if err != nil {
		log.Println("服务器发生错误：", err.Error())
		resp := &pb.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取用户失败",
		}
		return resp, nil
	}
	//	check 用户不存在
	if usr == nil {
		// log.Fatalln("用户不存在")
		resp := &pb.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		return resp, nil
	}
	//	check Avatar 头像
	//	TODO
	avatar := "Mio.jpg"
	//	check BackgroundImage 背景
	//	TODO
	backgroundImage := "kon.png"
	//	返回
	resp = &pb.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User: &pb.User{
			Id:              int64(usr.ID),
			Name:            usr.UserName,
			FollowCount:     int64(usr.FollowCount),
			FollowerCount:   int64(usr.FollowerCount),
			IsFollow:        usrID == int64(usr.ID), //	!TODO
			Avatar:          avatar,
			BackgroundImage: backgroundImage,
			Signature:       usr.Signature,
			TotalFavorited:  int64(usr.TotalFavorited),
			WorkCount:       int64(usr.WorkCount),
			FavoriteCount:   int64(usr.FavoriteCount),
		},
	}
	return resp, nil
}

// Login
func (this *UserServiceImpl) Login(ctx context.Context, req *pb.UserLoginRequest) (resp *pb.UserLoginResponse, err error) {
	/*
		username
		password
	*/
	//	Get user
	usr, err := db.GetUserByName(ctx, req.Username)
	//	check get error
	if err != nil {
		log.Println("获取用户失败：", err.Error())
		resp = &pb.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取用户失败",
		}
		return resp, err
	}
	//	check no such user
	if usr == nil {
		log.Println("用户名不存在")
		resp = &pb.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名不存在",
		}
		return resp, nil
	}
	//	check password
	if req.Password != usr.Password {
		log.Println("用户名或密码错误")
		resp = &pb.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "用户名或密码错误",
		}
		return resp, nil
	}
	//	jwt
	claims := jwt.CustomClaims{
		Id: int64(usr.ID),
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token, err := Jwt.GenToken(claims)
	if err != nil {
		log.Println("token创建失败：", err.Error())
		resp = &pb.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：token创建失败",
		}
		return resp, nil
	}
	//	return
	resp = &pb.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      token,
	}
	return resp, nil
}

// Register
func (this *UserServiceImpl) Register(ctx context.Context, req *pb.UserRegisterRequest) (resp *pb.UserRegisterResponse, err error) {
	//	check error
	// usr, err := db.GetUserByName(ctx, req.Username)
	// if err != nil {
	// 	log.Fatalln("检查用户失败：%v", err.Error())
	// 	resp = &pb.UserRegisterResponse{
	// 		StatusCode: -1,
	// 		StatusMsg:  "服务器错误：检查用户失败",
	// 	}
	// 	return resp, err
	// }
	//	check username
	if fl, _ := db.CheckNameRepeat(ctx, req.Username); fl == true {
		log.Println("用户名已存在：", req.Username)
		resp = &pb.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "用户名已存在",
		}
		return resp, nil
	}
	//	new user
	avaList := []string{"Azusa", "Mio", "Ritsu", "Tsumugi", "Yui"}
	rand.Seed(time.Now().UnixMilli())
	usr := &db.User{
		UserName: req.Username,
		Password: req.Password,
		Avatar:   fmt.Sprintf("%s.jpg", avaList[rand.Intn(5)]),
	}
	err = db.CreateUser(ctx, usr)
	if err != nil {
		log.Println("注册失败：", err.Error())
		resp = &pb.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：注册失败",
		}
		return resp, err
	}
	//	jwt
	claims := jwt.CustomClaims{
		Id: int64(usr.ID),
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token, err := Jwt.GenToken(claims)
	if err != nil {
		log.Println("token创建失败：", err.Error())
		resp = &pb.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：token创建失败",
		}
		return resp, nil
	}
	//	return
	resp = &pb.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(usr.ID),
		Token:      token,
	}
	return resp, nil
}
