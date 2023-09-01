package db

import (
	"context"

	"github.com/Road-To-Byte/Road2TikTok/pkg/errno"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ============ FavoriteVideoRelation 用户视频点赞数据结构 ============
type FavoriteVideoRelation struct {
	Video   Video `gorm:"foreignkey:VideoID;" json:"video,omitempty"`
	VideoID uint  `gorm:"index:idx_videoid;not null" json:"video_id"`
	User    User  `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID  uint  `gorm:"index:idx_userid;not null" json:"user_id"`
}

// ============ FavoriteCommentRelation 用户评论点赞数据结构 ============
type FavoriteCommentRelation struct {
	Comment   Comment `gorm:"foreignkey:CommentID;" json:"comment,omitempty"`
	CommentID uint    `gorm:"column:comment_id;index:idx_commentid;not null" json:"video_id"`
	User      User    `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID    uint    `gorm:"column:user_id;index:idx_userid;not null" json:"user_id"`
}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

func (FavoriteCommentRelation) TableName() string {
	return "user_favorite_comments"
}

// 添加 视频点赞
func CreateVideoFavorite(ctx context.Context, userID int64, videoID int64, authorID int64) error {
	//	Transaction 里面是一个func 这里把操作搞成一个事务 然后在一个func里一起操作
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//1. 新增点赞数据
		err := tx.Create(&FavoriteVideoRelation{UserID: uint(userID), VideoID: uint(videoID)}).Error
		if err != nil {
			return err
		}

		//2.改变 video 表中的 favorite count
		res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响的数据条数不是1
			return errno.ErrDatabase
		}

		//3.改变 user 表中的 favorite count
		res = tx.Model(&User{}).Where("id = ?", userID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		//4.改变 user 表中的 total_favorited
		res = tx.Model(&User{}).Where("id = ?", authorID).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 根据 VideoID和UserID 返回 点赞关系
func GetFavoriteVideoRelationByUserVideoID(ctx context.Context, userID int64, videoID int64) (*FavoriteVideoRelation, error) {
	FavoriteVideoRelation := new(FavoriteVideoRelation)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).First(&FavoriteVideoRelation, "user_id = ? and video_id = ?", userID, videoID).Error; err == nil {
		return FavoriteVideoRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 删除 视频点赞
func DelFavoriteByUserVideoID(ctx context.Context, userID int64, videoID int64, authorID int64) error {
	//	Transaction 里面是一个func 这里把操作搞成一个事务 然后在一个func里一起操作
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		FavoriteVideoRelation := new(FavoriteVideoRelation)
		if err := tx.Where("user_id = ? and video_id = ?", userID, videoID).First(&FavoriteVideoRelation).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}

		//1. 删除点赞数据
		// 因为FavoriteVideoRelation中包含了gorm.Model所以拥有软删除能力
		// 而tx.Unscoped().Delete()将永久删除记录
		err := tx.Unscoped().Where("user_id = ? and video_id = ?", userID, videoID).Delete(&FavoriteVideoRelation).Error
		if err != nil {
			return err
		}

		//2.改变 video 表中的 favorite count
		res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响数据条数不是1
			return errno.ErrDatabase
		}

		//3.改变 user 表中的 favorite count
		res = tx.Model(&User{}).Where("id = ?", userID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		//4.改变 user 表中的 total_favorited
		res = tx.Model(&User{}).Where("id = ?", authorID).Update("total_favorited", gorm.Expr("total_favorited - ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 根据 UserID 返回 视频点赞列表
func GetFavoriteListByUserID(ctx context.Context, userID int64) ([]*FavoriteVideoRelation, error) {
	var FavoriteVideoRelationList []*FavoriteVideoRelation
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("user_id = ?", userID).Find(&FavoriteVideoRelationList).Error
	if err != nil {
		return nil, err
	}
	return FavoriteVideoRelationList, nil
}

// 返回视频点赞列表
func GetAllFavoriteList(ctx context.Context) ([]*FavoriteVideoRelation, error) {
	var FavoriteVideoRelationList []*FavoriteVideoRelation
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Find(&FavoriteVideoRelationList).Error; err != nil {
		return nil, err
	}
	return FavoriteVideoRelationList, nil
}

// 根据 UserID和CommentID 返回 评论点赞关系
func GetFavoriteCommentRelationByUserCommentID(ctx context.Context, userID int64, commentID int64) (*FavoriteCommentRelation, error) {
	FavoriteCommentRelation := new(FavoriteCommentRelation)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).First(&FavoriteCommentRelation, "user_id = ? and comment_id = ?", userID, commentID).Error; err == nil {
		return FavoriteCommentRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
