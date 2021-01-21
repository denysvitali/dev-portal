package v1

import (
	"github.com/denysvitali/dev-portal/pkg/server/app"
	"github.com/denysvitali/dev-portal/pkg/server/routes/api/v1/topics"
	"github.com/denysvitali/dev-portal/pkg/server/routes/api/v1/users"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.RouterGroup, app *app.App){
	topics.Setup(r.Group("/topics"), app)
	users.Setup(r.Group("/users"), app)
}