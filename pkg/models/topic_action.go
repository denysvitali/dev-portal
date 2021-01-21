package models

import "time"

type TopicAction struct {
	ID        uint   `gorm:"primaryKey,autoincrement"`
	UserID    uint   `gorm:"uniqueIndex:idx_topic_action"`
	TopicID   uint   `gorm:"uniqueIndex:idx_topic_action"`
	ActionID  uint   `gorm:"uniqueIndex:idx_topic_action"`
	Action    Action `gorm:"foreignKey:ActionID"`
	CreatedAt time.Time
}
