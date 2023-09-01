package db

import (
	"context"

	"github.com/Road-To-Byte/Road2TikTok/pkg/errno"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

/*
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}
*/

// ============ FollowRelation 用户关注数据结构 ============
type FollowRelation struct {
	gorm.Model
	User     User `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID   uint `gorm:"index:idx_userid;not null" json:"user_id"`
	ToUser   User `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID uint `gorm:"index:idx_userid;index:idx_userid_to;not null" json:"to_user_id"`
}

//	============ 数据库修改操作 ===========

func (FollowRelation) TableName() string {
	return "relations"
}

// 根据 UserID toUserID 返回 FollowRelation
func GetRelationByUserIDs(ctx context.Context, userID int64, toUserID int64) (*FollowRelation, error) {
	relation := new(FollowRelation)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("user_id=? AND to_user_id=?", userID, toUserID).First(&relation).Error; err == nil {
		return relation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 添加 FollowRelation
func CreateRelation(ctx context.Context, userID int64, toUserID int64) error {
	//	Transaction 里面是一个func 这里把操作搞成一个事务 然后在一个func里一起操作
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// FollowRelation
		err := tx.Create(&FollowRelation{UserID: uint(userID), ToUserID: uint(toUserID)}).Error
		if err != nil {
			return err
		}

		// user: following count
		res := tx.Model(&User{}).Where("id = ?", userID).Update("following_count", gorm.Expr("following_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// user: follower count
		res = tx.Model(&User{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 删除 FollowRelation
func DelRelationByUserIDs(ctx context.Context, userID int64, toUserID int64) error {
	//	Transaction 里面是一个func 这里把操作搞成一个事务 然后在一个func里一起操作
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//	删除之前 先确保存在这样的 FollowRelation
		relation := new(FollowRelation)
		if err := tx.Where("user_id = ? AND to_user_id=?", userID, toUserID).First(&relation).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}

		// FollowRelation
		// 因为Relation中包含了gorm.Model所以拥有软删除能力
		// 而tx.Unscoped().Delete()将永久删除记录
		// 所谓的软删除其实就是打上删除的tag而不抹除数据 便于恢复
		// 如果不需要恢复 那么显然应该直接删除
		err := tx.Unscoped().Delete(&relation).Error
		//err := tx.Delete(&relation).Error	//软删除
		if err != nil {
			return err
		}

		// user: following count
		res := tx.Model(&User{}).Where("id = ?", userID).Update("following_count", gorm.Expr("following_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		// user: follower count
		res = tx.Model(&User{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		return nil
	})
	return err
}

// 根据 UserID 返回 follow 对应的 []FollowRelation
func GetFollowingListByUserID(ctx context.Context, userID int64) ([]*FollowRelation, error) {
	var RelationList []*FollowRelation
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("user_id = ?", userID).Find(&RelationList).Error
	if err != nil {
		return nil, err
	}
	return RelationList, nil
}

// 根据 UserID 返回 follower 对应的 []FollowRelation
func GetFollowerListByUserID(ctx context.Context, toUserID int64) ([]*FollowRelation, error) {
	var RelationList []*FollowRelation
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("to_user_id = ?", toUserID).Find(&RelationList).Error
	if err != nil {
		return nil, err
	}
	return RelationList, nil
}

// friend = follow&&follower
// 根据 UserID 返回 follow&&follower 对应的 []FollowRelation
func GetFriendList(ctx context.Context, userID int64) ([]*FollowRelation, error) {
	var FriendList []*FollowRelation
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Raw("SELECT user_id, to_user_id, created_at FROM relations WHERE user_id = ? AND to_user_id IN (SELECT user_id FROM relations r WHERE r.to_user_id = relations.user_id)", userID).Scan(&FriendList).Error
	if err != nil {
		return nil, err
	}
	return FriendList, nil
}
