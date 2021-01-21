package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	Username  string    `gorm:"index" json:"username"`
	GivenName string    `json:"given_name"`
	LastName  string    `json:"last_name"`
	Admin     bool      `gorm:"index" json:"admin"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"-"`
	Deleted   bool      `gorm:"index" json:"-"`

	UserDetails UserDetails `gorm:"foreignKey:ID" json:"details,omitempty"`
}
