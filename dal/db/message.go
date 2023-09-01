/*
 * @Autor: violet apricity ( Zhuangpx )
 * @Date: 2023-09-01 16:37:00
 * @LastEditors: violet apricity ( Zhuangpx )
 * @LastEditTime: 2023-09-01 17:09:58
 * @FilePath: \Road2TikTok\dal\db\message.go
 * @Description:  Zhuangpx : Violet && Apricity:/ The warmth of the sun in the winter /
 */
package db

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ============ Message 聊天消息数据结构 ============
type Message struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_time"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUser   User           `gorm:"foreignkey:FromUserID;" json:"from_user,omitempty"`
	FromUserID uint           `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	ToUser     User           `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID   uint           `gorm:"index:idx_userid_from;index:idx_userid_to;not null" json:"to_user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
}

func (Message) TableName() string {
	return "messages"
}

// 根据两个 UserID 返回二者的 Message
func GetMessagesByUserIDs(ctx context.Context, userID int64, toUserID int64, lastTimestamp int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND created_at > ?",
		userID, toUserID, toUserID, userID, time.UnixMilli(lastTimestamp).Format("2006-01-02 15:04:05.000"),
	).Order("created_at ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// 根据两个 UserID 返回user->touser的 Message
func GetMessagesByUserToUser(ctx context.Context, userID int64, toUserID int64, lastTimestamp int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Where("from_user_id = ? AND to_user_id = ? AND created_at > ?",
		userID, toUserID, time.UnixMilli(lastTimestamp).Format("2006-01-02 15:04:05.000"),
	).Order("created_at ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// 添加聊天信息
func CreateMessagesByList(ctx context.Context, messages []*Message) error {
	err := GetDB().Clauses(dbresolver.Write).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(messages).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 根据两个 UserID 返回 MessageID
func GetMessageIDsByUserIDs(ctx context.Context, userID int64, toUserID int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Select("id").Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userID, toUserID, toUserID, userID).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// 根据 MessageID 返回 Message
func GetMessageByID(ctx context.Context, messageID int64) (*Message, error) {
	res := new(Message)
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Select("id, from_user_id, to_user_id, content, created_at").Where("id = ?", messageID).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// 根据 Friend 返回最新 Message
func GetFriendLatestMessage(ctx context.Context, userID int64, toUserID int64) (*Message, error) {
	var res *Message
	if err := GetDB().Clauses(dbresolver.Read).WithContext(ctx).Select("id, from_user_id, to_user_id, content, created_at").Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userID, toUserID, toUserID, userID).Order("created_at DESC").Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
