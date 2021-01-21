package models

type UserDetails struct {
	ID         uint   `gorm:"primaryKey,autoincrement" json:"-"`
	Department string `gorm:"index" json:"department,omitempty"`
	Email      string `json:"email,omitempty"`
}
