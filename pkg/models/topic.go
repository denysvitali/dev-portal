package models

import (
	"time"
)

type Topic struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AuthorID      uint      `json:"-"`
	Author        User      `json:"author"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	CreatedAt     time.Time `json:"created_at"`
	Upvotes       uint      `gorm:"-" json:"upvotes"`
	Downvotes     uint      `gorm:"-" json:"downvotes"`
	TopicActions  []Action  `gorm:"many2many:topic_actions;" json:"-"`
	Comments      []Comment `gorm:"foreignKey:TopicID" json:"comments,omitempty"`
	CommentsCount uint      `gorm:"-" json:"comments_count"`
}

type Action struct {
	ID   uint `gorm:"primaryKey,autoIncrement"`
	Name string
}

type TopicAction struct {
	UserID    uint `gorm:"primaryKey"`
	TopicID   uint `gorm:"primaryKey"`
	ActionID  uint
	Action    Action `gorm:"foreignKey:ActionID"`
	CreatedAt time.Time
}
