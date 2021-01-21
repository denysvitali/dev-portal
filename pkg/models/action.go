package models

type Action struct {
	ID   uint `gorm:"primaryKey,autoIncrement"`
	Name string
}