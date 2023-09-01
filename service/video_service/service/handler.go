/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-31 21:54:57
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-09-02 03:43:06
 * @FilePath: \Road2TikTok\service\video_service\service\handler.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */

package video_service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Road-To-Byte/Road2TikTok/dal/db"
	"github.com/Road-To-Byte/Road2TikTok/pkg/viper"
	"github.com/Road-To-Byte/Road2TikTok/service/video_service/pb"
)

type VideoServiceImpl struct{}

const videoCntLimit = 30

// Feed
func (this *VideoServiceImpl) Feed(ctx context.Context, req *pb.FeedRequest) (resp *pb.FeedResponse, err error) {
	nextTime := time.Now().UnixMilli()
	var userID int64 = -1
	//	Check token
	if req.Token != "" {
		claims, err := Jwt.ParseToken(req.Token)
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			res := &pb.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：token 解析错误",
			}
			return res, nil
		}
		userID = claims.Id
	}
	// dal取video数据
	videos, err := db.GetLatestVideos(ctx, videoCntLimit, &req.LatestTime)
	if err != nil {
		log.Println("服务器发生错误：", err.Error())
		res := &pb.FeedResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取视频失败",
		}
		return res, nil
	}
	//	处理从dal取出的video数据
	videoList := make([]*pb.Video, 0)
	for _, vid := range videos {
		//	作者
		author, err := db.GetUserByUserID(ctx, int64(vid.AuthorID))
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			return nil, err
		}
		//	与作者的关系
		relation, err := db.GetRelationByUserIDs(ctx, userID, int64(author.ID))
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			res := &pb.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：获取视频作者关系失败",
			}
			return res, nil
		}
		//	与视频点赞关系
		favorite, err := db.GetFavoriteVideoRelationByUserVideoID(ctx, userID, int64(vid.ID))
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			res := &pb.FeedResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：获取视频点赞关系失败",
			}
			return res, nil
		}
		//	视频url
		playUrl := vid.PlayURL
		//	封面url
		coverUrl := vid.CoverURL
		//	作者头像url
		avatarUrl := author.Avatar
		//	作者主页背景url
		backgroundUrl := author.BackgroundImage
		//	视频列表
		videoList = append(videoList, &pb.Video{
			Id: int64(vid.ID),
			Author: &pb.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowCount:     int64(author.FollowCount),
				FollowerCount:   int64(author.FollowerCount),
				IsFollow:        relation != nil,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(vid.FavoriteCount),
			CommentCount:  int64(vid.CommentCount),
			IsFavorite:    favorite != nil,
			Title:         vid.Title,
		})
	}
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}
	res := &pb.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
		NextTime:   nextTime,
	}
	return res, nil
}

// PublishAction
func (this *VideoServiceImpl) PublishAction(ctx context.Context, req *pb.PublishActionRequest) (resp *pb.PublishActionResponse, err error) {
	//	解析token
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		log.Println("服务器发生错误：", err.Error())
		res := &pb.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：token 解析错误",
		}
		return res, nil
	}
	//	获取用户id
	userID := claims.Id
	//	检查标题
	if len(req.Title) == 0 || len(req.Title) > 32 {
		log.Println("标题不能为空且不能超过32个字符：", len(req.Title))
		res := &pb.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "客户端错误：标题不能为空且不能超过32个字符",
		}
		return res, nil
	}
	//	限制文件上传大小
	maxSize := viper.Init("video").Viper.GetInt("video.maxSizeLimit")
	size := len(req.Data)
	//	video.maxSizeLimit 的 1000*1000 倍 才是实际文件大小
	if size > maxSize*1000*1000 {
		log.Println("视频文件过大：", size)
		res := &pb.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  fmt.Sprintf("客户端错误：该视频文件大小(%dMB)超出限制(%dMB)", size, maxSize),
		}
		return res, nil
	}
	//	视频信息：视频标题 封面标题
	createTimestamp := time.Now().UnixMilli()
	videoTitle := fmt.Sprintf("%d_%s_%d.mp4", userID, req.Title, createTimestamp)
	coverTitle := fmt.Sprintf("%d_%s_%d.png", userID, req.Title, createTimestamp)

	//	更新db
	v := &db.Video{
		Title:    req.Title,
		PlayURL:  videoTitle,
		CoverURL: coverTitle,
		AuthorID: uint(userID),
	}
	err = db.InsertVideo(ctx, v)
	if err != nil {
		log.Println("服务器发生错误：", err.Error())
		res := &pb.PublishActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：视频发布失败",
		}
		return res, nil
	}

	//	开个协程 更新video 如果发布失败 就把记录删除
	go func() {
		err := VideoPublish(req.Data, videoTitle, coverTitle)
		if err != nil {
			//	发生错误，则删除插入的记录
			e := db.DelVideoByID(ctx, int64(v.ID), userID)
			if e != nil {
				log.Println("视频记录删除失败：", err.Error())
			}
		}
	}()

	//	返回响应
	res := &pb.PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "视频发布成功，等待审核",
	}
	return res, nil
}

// PublishList
func (this *VideoServiceImpl) PublishList(ctx context.Context, req *pb.PublishListRequest) (resp *pb.PublishListResponse, err error) {
	//	获取用户ID
	userID := req.UserId
	//	从dal获取用户视频发布列表
	videoList, err := db.GetVideosByAuthorID(ctx, userID)
	if err != nil {
		log.Println("服务器发生错误：", err.Error())
		res := &pb.PublishListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：发布列表获取失败",
		}
		return res, nil
	}
	//	处理从dal取出的数据
	videos := make([]*pb.Video, 0)
	for _, vid := range videoList {
		//	作者
		author, err := db.GetUserByUserID(ctx, int64(vid.AuthorID))
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			res := &pb.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：发布列表作者获取失败",
			}
			return res, nil
		}
		//	作者粉丝
		follow, err := db.GetRelationByUserIDs(ctx, userID, int64(author.ID))
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			res := &pb.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：发布列表作者粉丝获取失败",
			}
			return res, nil
		}
		//	视频点赞
		favorite, err := db.GetFavoriteVideoRelationByUserVideoID(ctx, userID, int64(vid.ID))
		if err != nil {
			log.Println("服务器发生错误：", err.Error())
			res := &pb.PublishListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：点赞信息获取失败",
			}
			return res, nil
		}
		//	视频url
		playUrl := vid.PlayURL
		//	封面url
		coverUrl := vid.CoverURL
		//	作者头像url
		avatarUrl := author.Avatar
		//	作者主页背景url
		backgroundUrl := author.BackgroundImage
		//	视频发布列表
		videos = append(videos, &pb.Video{
			Id: int64(vid.ID),
			Author: &pb.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowCount:     int64(author.FollowCount),
				FollowerCount:   int64(author.FollowerCount),
				IsFollow:        follow != nil,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(vid.FavoriteCount),
			CommentCount:  int64(vid.CommentCount),
			IsFavorite:    favorite != nil,
			Title:         vid.Title,
		})
	}
	//	返回
	res := &pb.PublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videos,
	}
	return res, nil
}
