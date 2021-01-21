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
	CreatedAt     time.Time `json:"createdAt"`
	Likes         uint      `gorm:"-" json:"likes"`
	Liked         bool      `gorm:"-" json:"liked"`
	TopicActions  []Action  `gorm:"many2many:topic_actions;" json:"-"`
	Comments      []Comment `gorm:"foreignKey:TopicID" json:"comments,omitempty"`
	CommentsCount uint      `gorm:"-" json:"commentsCount"`
}
