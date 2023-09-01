package db

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Video struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time `gorm:"not null;index:update_idx;column:update_time;"`
	Author        User      `gorm:"foreignkey:AuthorID"`
	AuthorID      uint      `gorm:"not null;index:AuthorID_idx"`
	PlayURL       string    `gorm:"not null;type:varchar(255)"`
	CoverURL      string    `gorm:"type:varchar(255)"`
	FavoriteCount int       `gorm:"not null;default:0"`
	CommentCount  int       `gorm:"not null;default:0"`
	Title         string    `gorm:"not null;type:varchar(255)"`
}

func (Video) TableName() string {
	return "videos"
}

func GetVideoByID(ctx context.Context, videoID int64) (*Video, error) {
	video := new(Video)
	// GetDB().Raw("SELECT * FROM videos WHERE id = ?", videoID).Scan(&ret)
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("id = ?", videoID).First(&video).Error
	if err == nil {
		return video, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return nil, err
	}
}

// GetVideoListByIDs may never be used
func GetVideoListByIDs(ctx context.Context, videoIDs []int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	num := len(videoIDs)
	if num == 0 {
		return videos, nil
	}
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("id in ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, err
	}

	return videos, nil
}

// GetLatestVideos returns <limit> videos according to <latestTime>
func GetLatestVideos(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	if latestTime == nil || *latestTime == 0 {
		cur := time.Now().UnixMilli()
		*latestTime = cur
	}
	// SELECT * FROM videos WHERE update_time < <latest_time> LIMIT <limit> ORDER BY update_time desc
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Limit(limit).Order("update_time desc").
		Where("update_time < ?", time.UnixMilli(*latestTime)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
