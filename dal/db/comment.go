package db

import (
	"context"
	"time"

	"github.com/Road-To-Byte/Road2TikTok/pkg/errno"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ============ Comment 用户评论数据结构 ============
type Comment struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_date"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Video      Video          `gorm:"foreignkey:VideoID" json:"video,omitempty"`
	VideoID    uint           `gorm:"index:idx_videoid;not null" json:"video_id"`
	User       User           `gorm:"foreignkey:UserID" json:"user,omitempty"`
	UserID     uint           `gorm:"index:idx_userid;not null" json:"user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
	LikeCount  uint           `gorm:"column:like_count;default:0;not null" json:"like_count,omitempty"`
	TeaseCount uint           `gorm:"column:tease_count;default:0;not null" json:"tease_count,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

// 添加评论
func CreateComment(ctx context.Context, comment *Comment) error {
	//	Transaction 里面是一个func 这里把操作搞成一个事务 然后在一个func里一起操作
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 新增评论数据
		err := tx.Create(comment).Error
		if err != nil {
			return err
		}

		// 2.对 Video 表中的评论数+1
		res := tx.Model(&Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响的数据条数不是1
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 删除评论
func DelCommentByID(ctx context.Context, commentID int64, vid int64) error {
	//	Transaction 里面是一个func 这里把操作搞成一个事务 然后在一个func里一起操作
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//	删除之前 先确保存在这样的 Comment
		comment := new(Comment)
		if err := tx.First(&comment, commentID).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}

		// 1. 删除评论数据
		// 这里使用的实际上是软删除
		err := tx.Where("id = ?", commentID).Delete(&Comment{}).Error
		if err != nil {
			return err
		}

		// 2.改变 video 表中的 comment count
		res := tx.Model(&Video{}).Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响的数据条数不是1
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

//	根据 VideoID 返回 Comment
func GetVideoCommentListByVideoID(ctx context.Context, videoID int64) ([]*Comment, error) {
	var comments []*Comment
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Model(&Comment{}).Where(&Comment{VideoID: uint(videoID)}).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

//	根据 CommentID 返回 Comment
func GetCommentByCommentID(ctx context.Context, commentID int64) (*Comment, error) {
	comment := new(Comment)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("id = ?", commentID).First(&comment).Error; err == nil {
		return comment, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
