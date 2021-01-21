package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	Username  string    `gorm:"index" json:"username"`
	GivenName string    `json:"givenName"`
	LastName  string    `json:"lastName"`
	Admin     bool      `gorm:"index" json:"admin"`
	CreatedAt time.Time `json:"createdAt"`
	LastLogin time.Time `json:"-"`
	Deleted   bool      `gorm:"index" json:"-"`

	UserDetails UserDetails `gorm:"foreignKey:ID" json:"details,omitempty"`
}
