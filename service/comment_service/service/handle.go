/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-09-02 15:45:38
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-09-02 17:43:07
 * @FilePath: \Road2TikTok\service\comment_service\service\handle.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package comment_service

import (
	"context"
	"log"

	"github.com/Road-To-Byte/Road2TikTok/dal/db"
	"github.com/Road-To-Byte/Road2TikTok/service/comment_service/pb"
	"gorm.io/gorm"
)

type CommentServiceImpl struct{}

// CommentAction
func (this *CommentServiceImpl) CommentAction(ctx context.Context, req *pb.CommentActionRequest) (resp *pb.CommentActionResponse, err error) {
	//	解析token
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		log.Println("服务器发生错误，解析token失败：", err.Error())
		res := &pb.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：token 解析错误",
		}
		return res, nil
	}
	//	获取用户id 视频 评论行为
	userID := claims.Id
	actionType := req.ActionType
	vid, err := db.GetVideoByID(ctx, req.VideoId)
	if vid == nil {
		log.Println("服务器发生错误，视频不存在：", err.Error())
		res := &pb.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：视频不存在",
		}
		return res, nil
	}
	//	1-添加评论 2-删除评论
	if actionType == 1 {
		//	new comment
		cmt := &db.Comment{
			VideoID: uint(req.VideoId),
			UserID:  uint(userID),
			Content: req.CommentText,
		}
		//	更新db
		err := db.CreateComment(ctx, cmt)
		if err != nil {
			log.Println("服务器发生错误，添加评论失败：", err.Error())
			res := &pb.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：发布评论失败",
			}
			return res, nil
		}
	} else if actionType == 2 {
		//	删除评论前 确保该评论是否发布自该用户，或该评论在该用户所发布的视频下
		//	get comment
		cmt, err := db.GetCommentByCommentID(ctx, req.CommentId)
		if err != nil {
			log.Println("服务器发生错误，查找评论失败：", err.Error())
			res := &pb.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：查找评论失败",
			}
			return res, nil
		}
		if cmt == nil {
			//	评论不存在，无法删除
			log.Println("评论删除失败，找不到该评论ID：", req.CommentId)
			res := &pb.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "评论删除失败，找不到该评论ID",
			}
			return res, nil
		} else {
			//	视频
			v, err := db.GetVideoByID(ctx, int64(cmt.VideoID))
			if err != nil {
				log.Println("服务器发生错误，评论删除失败：", err.Error())
				res := &pb.CommentActionResponse{
					StatusCode: -1,
					StatusMsg:  "评论删除失败：服务器内部错误",
				}
				return res, nil
			}
			//	若删除评论的用户不是发布评论的用户且该用户不是视频作者
			if userID != int64(cmt.UserID) && userID != int64(v.AuthorID) {
				log.Println("评论删除失败，没有权限：", cmt.UserID)
				res := &pb.CommentActionResponse{
					StatusCode: -1,
					StatusMsg:  "评论删除失败：没有权限",
				}
				return res, nil
			}
		}
		err = db.DelCommentByID(ctx, req.CommentId, req.VideoId)
		if err != nil {
			log.Println("服务器发生错误，评论删除失败：", err.Error())
			res := &pb.CommentActionResponse{
				StatusCode: -1,
				StatusMsg:  "评论删除失败：服务器内部错误",
			}
			return res, nil
		}
	} else {
		res := &pb.CommentActionResponse{
			StatusCode: -1,
			StatusMsg:  "action_type 非法",
		}
		return res, nil
	}
	res := &pb.CommentActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	return res, nil
}

// CommentList
func (this *CommentServiceImpl) CommentList(ctx context.Context, req *pb.CommentListRequest) (resp *pb.CommentListResponse, err error) {
	var userID int64 = -1
	//	解析token
	if req.Token != "" {
		claims, err := Jwt.ParseToken(req.Token)
		if err != nil {
			log.Println("服务器发生错误，解析token失败：", err.Error())
			res := &pb.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "token 解析错误",
			}
			return res, nil
		}
		userID = claims.Id
	}

	//	dal取数据
	results, err := db.GetVideoCommentListByVideoID(ctx, req.VideoId)
	if err != nil {
		log.Println("服务器发生错误，获取评论列表失败：", err.Error())
		res := &pb.CommentListResponse{
			StatusCode: -1,
			StatusMsg:  "服务器错误：获取评论列表失败",
		}
		return res, nil
	}
	//	处理dal数据
	comments := make([]*pb.Comment, 0)
	for _, r := range results {
		//	用户
		u, err := db.GetUserByUserID(ctx, int64(r.UserID))
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Println("服务器发生错误，获取用户错误：", err.Error())
			res := &pb.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：获取评论用户失败",
			}
			return res, nil
		}
		//	用户关系
		_, err = db.GetRelationByUserIDs(ctx, userID, int64(u.ID))
		if err != nil {
			log.Println("服务器发生错误，获取用户关系错误：", err.Error())
			res := &pb.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器错误：获取用户关系失败",
			}
			return res, nil
		}
		//	头像
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio获取头像失败：%v", err.Error())
			res := &pb.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误：评论列表用户头像获取失败",
			}
			return res, nil
		}
		//	背景
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio获取背景图链接失败：%v", err.Error())
			res := &pb.CommentListResponse{
				StatusCode: -1,
				StatusMsg:  "服务器内部错误：评论列表用户背景图获取失败",
			}
			return res, nil
		}
		usr := &pb.User{
			Id:              userID,
			Name:            u.UserName,
			FollowCount:     int64(u.FollowCount),
			FollowerCount:   int64(u.FollowerCount),
			IsFollow:        err != gorm.ErrRecordNotFound,
			Avatar:          avatar,
			BackgroundImage: backgroundUrl,
			Signature:       u.Signature,
			TotalFavorited:  int64(u.TotalFavorited),
			WorkCount:       int64(u.WorkCount),
			FavoriteCount:   int64(u.FavoriteCount),
		}
		comments = append(comments, &pb.Comment{
			Id:         int64(r.ID),
			User:       usr,
			Content:    r.Content,
			CreateDate: r.CreatedAt.Format("2006-01-02"),
			LikeCount:  int64(r.LikeCount),
			TeaseCount: int64(r.TeaseCount),
		})
	}

	res := &pb.CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: comments,
	}
	return res, nil
}
