/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-21 02:01:49
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-25 15:39:41
 * @FilePath: \Road2TikTok\service\user_service\service\handler.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package user_service

import (
	"context"
	"log"

	"github.com/Road-To-Byte/Road2TikTok/dal/db"
	"github.com/Road-To-Byte/Road2TikTok/service/user_service/pb"
)

type UserServiceImpl struct{}

//	UserInfo
func (this *UserServiceImpl) UserInfo(ctx context.Context, req *pb.UserInfoRequest) (resp *pb.UserInfoResponse, err error) {
	usrID := req.UserId
	//	dal取出user
	usr, err := db.GetUserByUserID(ctx, usrID)
	//	check err
	if err != nil {
		log.Fatalln("服务器发送错误：", err.Error())
		resp := &pb.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "服务器内部错误：获取用户失败",
		}
		return resp, nil
	}
	//	check 用户不存在
	if usr == nil {
		log.Fatalln("用户不存在：%v", err.Error())
		resp := &pb.UserInfoResponse{
			StatusCode: -1,
			StatusMsg:  "用户不存在",
		}
		return resp, nil
	}
	//	check Avatar 头像
	//	TODO
	//	check BackgroundImage 背景
	//	TODO
	//	返回
	resp = &pb.UserInfoResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User: &pb.User{
			Id:            int64(usr.ID),
			Name:          usr.UserName,
			FollowCount:   int64(usr.FollowCount),
			FollowerCount: int64(usr.FollowerCount),
			IsFollow:      usrID == int64(usr.ID), //	!TODO
			// Avatar: avatar,
			// BackgroundImage: backgroundImage,
			Signature:      usr.Signature,
			TotalFavorited: int64(usr.TotalFavorited),
			WorkCount:      int64(usr.WorkCount),
			FavoriteCount:  int64(usr.FavoriteCount),
		},
	}
	return resp, nil
}
