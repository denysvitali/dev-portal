package api

import (
	"github.com/denysvitali/dev-portal/pkg/server/app"
	v1 "github.com/denysvitali/dev-portal/pkg/server/routes/api/v1"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.RouterGroup, app *app.App){
	v1.Setup(r.Group("/v1"), app)
}