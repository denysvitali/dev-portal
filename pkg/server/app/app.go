package app

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	Db *gorm.DB
	Log *logrus.Logger
}