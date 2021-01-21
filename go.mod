module github.com/denysvitali/dev-portal

go 1.15

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/sirupsen/logrus v1.7.0
	gorm.io/driver/postgres v1.0.6
	gorm.io/gorm v1.20.11
)

replace gorm.io/gorm => ../../go-gorm/gorm
