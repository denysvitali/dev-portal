package topics

import (
	"github.com/denysvitali/dev-portal/pkg/models"
	"github.com/denysvitali/dev-portal/pkg/server/app"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
)

func Setup(r *gin.RouterGroup, app *app.App) {
	r.GET("/", getTopics(app))
	r.GET("/:id", getTopic(app))
}

func getTopic(a *app.App) gin.HandlerFunc {
	return func(context *gin.Context) {
		topicInput := context.Param("id")
		topicId, err := strconv.Atoi(topicInput)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "you need to provide a valid topic ID"})
			return
		}
		var topic models.Topic
		if err := a.Db.Preload("Author").Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.vote DESC, comments.created_at DESC")
		}).Preload("Comments.Author").First(&topic, topicId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				context.JSON(http.StatusNotFound, gin.H{"error": "topic not found"})
				return
			}
			a.Log.Errorf("unable to fetch topic %d: %v", topicId, err)
			context.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		err = enrichTopic(&topic, a)
		if err != nil {
			a.Log.Errorf("unable to enrich topic: %v", err)
			context.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		context.JSON(http.StatusOK, topic)
	}
}

func getLikesByTopic(id uint, app *app.App) (uint, error) {
	var upvotes int64 = 0
	err := app.Db.Table("topic_actions").Joins("JOIN actions ON topic_actions.action_id = actions.id").Group("topic_id").Group("actions.id").Where("topic_id=? AND Actions.name=?", id, "like").Count(&upvotes).Error
	return uint(upvotes), err
}

func getDownvotesByTopic(id uint, app *app.App) (uint, error) {
	var downvotes int64 = 0
	err := app.Db.Table("topic_actions").Joins("JOIN actions ON topic_actions.action_id = actions.id").Group("topic_id").Group("actions.id").Where("topic_id=? AND Actions.name=?", id, "downvote").Count(&downvotes).Error
	return uint(downvotes), err
}

func getCommentsCountByTopic(id uint, app *app.App) (uint, error) {
	var comments int64 = 0
	err := app.Db.Table("comments").Where("topic_id=?", id).Count(&comments).Error
	return uint(comments), err
}


func getTopics(app *app.App) func(context *gin.Context) {
	return func(context *gin.Context) {
		var topics []*models.Topic
		dbPreload := app.Db.Preload("Author")
		if err := dbPreload.Find(&topics).Error; err != nil {
			app.Log.Errorf("unable to get topics: %v", err)
		}

		for _, t := range topics {
			err := enrichTopic(t, app)
			if err != nil {
				app.Log.Errorf("unable to enrich topic: %v", err)
				context.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}
		context.JSON(http.StatusOK, topics)
	}
}

func enrichTopic(t *models.Topic, app *app.App) error {
	var err error
	t.Likes, err = getLikesByTopic(t.ID, app)
	if err != nil {
		return err
	}
	t.Liked = rand.Intn(2) == 1 // TODO: change to real value when we implement sessions
	t.CommentsCount, err = getCommentsCountByTopic(t.ID, app)
	
	return nil
}