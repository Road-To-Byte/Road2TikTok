/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-09-02 17:57:39
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-09-02 18:42:25
 * @FilePath: \Road2TikTok\service\favorite_service\service\handle.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package favorite_service

import (
	"context"
	"log"

	"github.com/Road-To-Byte/Road2TikTok/api_gateway/rpc/pb"
	"github.com/Road-To-Byte/Road2TikTok/dal/db"
	"github.com/Road-To-Byte/Road2TikTok/pkg/minio"
)

type FavoriteServiceImpl struct{}

func (this *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *pb.FavoriteActionRequest) (resp *pb.FavoriteActionResponse, err error) {
	//	解析token
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		log.Println("服务器发生错误，token解析失败：", err.Error())
		res := &pb.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：token 解析错误",
		}
		return res, nil
	}
	//	用户ID 视频ID 行为-点赞或者取消点赞
	userID := claims.Id
	actionType := req.ActionType
	vid, err := db.GetVideoByID(ctx, req.VideoId)
	if vid == nil {
		log.Println("服务器发生错误，视频不存在：", err.Error())
		res := &pb.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：视频不存在",
		}
		return res, nil
	}
	//	1-点赞 2取消点赞
	if actionType == 1 {
		//	更新db
		err = db.CreateVideoFavorite(ctx, userID, int64(vid.ID), int64(vid.AuthorID))
		if err != nil {
			log.Println("服务器发生错误，视频点赞失败：", err.Error())
			res := &pb.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：视频点赞失败",
			}
			return res, nil
		}
	} else if actionType == 2 {
		err = db.DelFavoriteByUserVideoID(ctx, userID, int64(vid.ID), int64(vid.AuthorID))
		if err != nil {
			log.Println("服务器发生错误，视频取消点赞失败：", err.Error())
			res := &pb.FavoriteActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：视频取消点赞失败",
			}
			return res, nil
		}
	} else {
		res := &pb.FavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "action_type 非法",
		}
		return res, nil
	}
	res := &pb.FavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return res, nil
}

// FavoriteList
func (this *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *pb.FavoriteListRequest) (resp *pb.FavoriteListResponse, err error) {
	var userID int64 = -1
	//	解析token
	if req.Token != "" {
		claims, err := Jwt.ParseToken(req.Token)
		if err != nil {
			log.Println("服务器发生错误，解析token失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "token 解析错误",
			}
			return res, nil
		}
		userID = claims.Id
	}
	//	dal取数据
	results, err := db.GetFavoriteListByUserID(ctx, userID)
	if err != nil {
		log.Println("服务器发生错误，获取喜欢列表失败：", err.Error())
		res := &pb.FavoriteListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取喜欢列表失败",
		}
		return res, nil
	}
	//	处理dal数据
	favoriteVideos := make([]*pb.Video, 0)
	for _, fav := range results {
		//	视频
		vid, err := db.GetVideoByID(ctx, int64(fav.VideoID))
		if err != nil {
			log.Println("服务器发生错误，获取喜欢视频失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：获取喜欢视频失败",
			}
			return res, nil
		}
		//	作者
		aut, err := db.GetUserByUserID(ctx, int64(vid.AuthorID))
		if err != nil {
			log.Println("服务器发生错误，获取喜欢视频作者失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：获取喜欢视频作者失败",
			}
			return res, nil
		}
		//	与作者的用户关系
		relation, err := db.GetRelationByUserIDs(ctx, userID, int64(aut.ID))
		if err != nil {
			log.Println("服务器发生错误，获取作者关系失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误，获取作者关系失败",
			}
			return res, nil
		}
		//	视频url
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, vid.PlayURL)
		if err != nil {
			log.Println("服务器发生错误，获取视频链接失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误，获取视频链接失败",
			}
			return res, nil
		}
		//	封面url
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, vid.CoverURL)
		if err != nil {
			log.Println("服务器发生错误，获取视频封面失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误，获取视频封面失败",
			}
			return res, nil
		}
		//	作者头像
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, aut.Avatar)
		if err != nil {
			log.Println("服务器发生错误，获取作者头像失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误，获取作者头像失败",
			}
			return res, nil
		}
		//	作者主页背景
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, aut.BackgroundImage)
		if err != nil {
			log.Println("服务器发生错误，获取作者背景失败：", err.Error())
			res := &pb.FavoriteListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误，获取作者背景失败",
			}
			return res, nil
		}
		favoriteVideos = append(favoriteVideos, &pb.Video{
			Id: int64(fav.VideoID),
			Author: &pb.User{
				Id:              int64(aut.ID),
				Name:            aut.UserName,
				FollowCount:     int64(aut.FollowCount),
				FollowerCount:   int64(aut.FollowerCount),
				IsFollow:        relation != nil,
				Avatar:          avatar,
				BackgroundImage: backgroundUrl,
				Signature:       aut.Signature,
				TotalFavorited:  int64(aut.TotalFavorited),
				WorkCount:       int64(aut.WorkCount),
				FavoriteCount:   int64(aut.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(vid.FavoriteCount),
			CommentCount:  int64(vid.CommentCount),
			IsFavorite:    true,
			Title:         vid.Title,
		})
	}
	res := &pb.FavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  favoriteVideos,
	}
	return res, nil
}
