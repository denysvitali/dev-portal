package models

import "time"

type Comment struct {
	ID              uint      `gorm:"primaryKey,autoincrement" json:"id"`
	Author          User      `gorm:"foreignKey:AuthorID" json:"author"`
	AuthorID        uint      `json:"-"`
	TopicID         uint      `json:"topic_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
	ParentComment   *Comment   `gorm:"foreignKey:ParentCommentId" json:"parent_comment"`
	ParentCommentId uint      `json:"parent_comment_id"`
	Votes           int       `json:"votes"`
}
