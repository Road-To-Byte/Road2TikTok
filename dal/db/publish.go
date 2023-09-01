package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// InsertVideo inserts video into "videos" table
func InsertVideo(ctx context.Context, video *Video) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// INSERT INTO videos VALUE ()
		err := tx.Create(video).Error
		if err != nil {
			return err
		}

		// UPDATE users SET work_count = work_count + 1 WHERE id = <video.AuthorID>
		err = tx.Model("users").Where("id = ?", video.AuthorID).
			Update("work_count = ?", gorm.Expr("work_count + ?", 1)).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetVideosByAuthorID returns a list of *Video according to author ID
func GetVideosByAuthorID(ctx context.Context, authorID int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	// SELECT * FROM videos WHERE author_id = <authorID>
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).
		Where("author_id = ?", authorID).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// DelVideoByID deletes video from table "videos" according to author ID and video ID
func DelVideoByID(ctx context.Context, authorID int64, videoID int64) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Unscoped().Delete(&Video{}, videoID).Error
		if err != nil {
			return err
		}
		err = tx.Model(&User{}).Where("id = ?", authorID).
			Update("work_count", gorm.Expr("work_count - ?", 1)).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
