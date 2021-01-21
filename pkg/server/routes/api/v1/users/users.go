package users

import (
	"github.com/denysvitali/dev-portal/pkg/models"
	"github.com/denysvitali/dev-portal/pkg/server/app"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
)

func Setup(r *gin.RouterGroup, app *app.App){
	r.GET("/", getUsers(app))
}

func getUsers(app *app.App) func(context *gin.Context) {
	return func(context *gin.Context) {
		var users []models.User
		if err := app.Db.Preload(clause.Associations).Find(&users).Error; err != nil {
			app.Log.Errorf("unable to get users: %v", err)
		}
		context.JSON(http.StatusOK, users)
	}
}