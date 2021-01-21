package server

import (
	"github.com/denysvitali/dev-portal/pkg/models"
	"github.com/denysvitali/dev-portal/pkg/server/app"
	"github.com/denysvitali/dev-portal/pkg/server/routes/api"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Server struct {
	listenAddr string
	log        *logrus.Logger
}

func New(listenAddr string, logger *logrus.Logger) Server {
	return Server{
		log:        logger,
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://127.0.0.1:8080"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	r.Use(cors.New(corsConfig))

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres"), &gorm.Config{
		Logger: newLogger,
	})
	
	if err != nil {
		s.log.Fatalf("unable to open db: %v", err)
	}
	s.setupORM(db)
	s.initData(db)

	app := app.App{
		Db:  db,
		Log: s.log,
	}

	setupRoutes(r, &app)
	err = r.Run(s.listenAddr)
	if err != nil {
		s.log.Errorf("server run error: %v", err)
	}
}

func (s *Server) migrateOrFail(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(model)
	if err != nil {
		s.log.Fatalf("unable to migrate: %v", err)
	}
}

func (s *Server) setupORM(db *gorm.DB) {
	s.migrateOrFail(db, &models.UserDetails{})
	s.migrateOrFail(db, &models.User{})
	s.migrateOrFail(db, &models.Topic{})
	s.migrateOrFail(db, &models.TopicAction{})
	s.migrateOrFail(db, &models.Action{})
	s.migrateOrFail(db, &models.Comment{})
	_ = db.SetupJoinTable(&models.Topic{}, "TopicActions", &models.TopicAction{})
}

func (s *Server) initData(db *gorm.DB) {
	var adminUser models.User
	if err := db.Where("username = ?", "admin").Find(&adminUser).Error; err != nil {
		s.log.Fatalf("unable to get admin user: %v", err)
	}

	emptyUser := models.User{}

	transaction := db.Debug().Begin()
	if adminUser == emptyUser {
		// Create admin
		adminUser = models.User{
			Username:  "admin",
			GivenName: "Site",
			LastName:  "Admin",
			Admin:     true,
			CreatedAt: time.Now(),
			Deleted:   false,
			UserDetails: models.UserDetails{
				Department: "A-B-C",
				Email:      "admin@login.dev",
			},
		}
		transaction.Create(&adminUser)
		transaction.Save(&adminUser)
 		transaction.Commit()
	}

	var topicsCount int64
	db.Table("topics").Count(&topicsCount)

	if topicsCount == 0 {
		db.Create(&models.Topic{
			Author: adminUser,
			Title:  "First Topic",
			Body:   "Hello world!",
		})
		db.Commit()
	}
}

func setupRoutes(r *gin.Engine, app *app.App) {
	api.Setup(r.Group("/api"), app)
}
