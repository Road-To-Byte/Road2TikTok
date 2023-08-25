/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-08-22 17:19:49
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-08-24 17:07:15
 * @FilePath: \Road2TikTok\dal\db\user.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package db

import (
	"context"

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

//	============ User 用户数据结构 ============
type User struct {
	gorm.Model
	UserName        string `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password        string `gorm:"type:varchar(256);not null" json:"password,omitempty"`
	FollowCount     uint   `gorm:"default:0;not null" json:"follow_count,omitempty"`
	FollowerCount   uint   `gorm:"default:0;not null" json:"follower_count,omitempty"`
	Avatar          string `gorm:"type:varchar(256)" json:"avatar,omitempty"`
	BackgroundImage string `gorm:"column:background_image;type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"`
	TotalFavorited  int64  `gorm:"default:0;not null" json:"total_favorited,omitempty"`
	WorkCount       int64  `gorm:"default:0;not null" json:"work_count,omitempty"`
	FavoriteCount   int64  `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	// FavoriteVideos  []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos,omitempty"`
	Signature string `gorm:"type:varchar(256)" json:"signature,omitempty"`
}

//	============ 数据库修改操作 ===========
//	在进行读写操作的时候 通过dbresolver 主数据库写 从数据库读 实现读写分离 减小主数据库的压力

//	table name
func (this *User) TableName() string {
	return "users"
}

//	根据 []IDs 返回 []Users
func GetUsersByUserIDs(ctx context.Context, userIDs []int64) ([]*User, error) {
	users := make([]*User, 0)
	//	[]IDs 为空
	if len(userIDs) == 0 {
		return users, nil
	}
	//	查找
	//	这里是[]IDs 所以是in的集合操作
	// err := GetDB().WithContext(ctx).Where("id in ?", userIDs).Find(&users).Error
	err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("id in ?", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

//	根据 ID 返回 User
func GetUserByUserID(ctx context.Context, userID int64) (*User, error) {
	user := new(User)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).First(&user, userID).Error; err == nil {
		return user, err
	} else if err == gorm.ErrRecordNotFound {
		//	err: not fount
		return nil, nil
	} else {
		//	err: other
		return nil, err
	}
}

//	根据 Name 返回 User
func GetUserByName(ctx context.Context, userName string) (*User, error) {
	user := new(User)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Select("id, user_name, password").Where("user_name = ?", userName).First(&user).Error; err == nil {
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

// 根据 Name 返回 Password
func GetPasswordByUsername(ctx context.Context, userName string) (*User, error) {
	user := new(User)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Select("password").Where("user_name = ?", userName).First(&user).Error; err == nil {
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

//	添加 []Users
func CreateUsers(ctx context.Context, users []*User) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(users).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

//	添加 User
func CreateUser(ctx context.Context, user *User) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
