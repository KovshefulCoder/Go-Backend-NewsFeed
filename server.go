package main

import (
	"NewsFeed/pkg/api"
	"NewsFeed/pkg/config"
	"NewsFeed/pkg/services"
	"NewsFeed/pkg/services/auth"
	"NewsFeed/pkg/services/feed"
	"NewsFeed/pkg/services/user"
	"NewsFeed/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func setRoutes(r *gin.Engine, coreApi *api.Api) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	general := r.Group("/")
	general.POST("/auth/signup_number", func(c *gin.Context) {
		coreApi.AuthApi.SignUpNumber(c, c.Request)
	})
	general.POST("/auth/refresh", func(c *gin.Context) {
		coreApi.AuthApi.RefreshToken(c, c.Request)
	})

	authGroup := r.Group("/auth").Use(coreApi.AuthApi.AuthMW())
	authGroup.POST("/signin_number", func(c *gin.Context) {
		coreApi.AuthApi.SignInNumber(c, c.Request)
	})
	authGroup.POST("/code", func(c *gin.Context) {
		coreApi.AuthApi.AuthCode(c, c.Request)
	})
	authGroup.POST("/2fa", func(c *gin.Context) {
		coreApi.AuthApi.Auth2FA(c, c.Request)
	})
	authGroup.GET("/check_client", func(c *gin.Context) {
		coreApi.AuthApi.CheckClient(c, c.Request)
	})

	channelsGroup := r.Group("/channels").Use(coreApi.AuthApi.AuthMW())
	channelsGroup.GET("/get_feed", func(c *gin.Context) {
		coreApi.FeedAPI.GetFeedChannels(c, c.Request)
	})
	channelsGroup.POST("/add_feed", func(c *gin.Context) {
		coreApi.FeedAPI.AddFeedChannel(c, c.Request)
	})
	channelsGroup.DELETE("/remove_feed", func(c *gin.Context) {
		coreApi.FeedAPI.RemoveFeedChannel(c, c.Request)
	})

	userGroup := r.Group("/user").Use(coreApi.AuthApi.AuthMW())
	userGroup.GET("/get_subs", func(c *gin.Context) {
		coreApi.UserApi.GetAllPublicSubsChannels(c, c.Request)
	})
	userGroup.GET("/form_feed", func(c *gin.Context) {
		coreApi.UserApi.FormFeed(c, c.Request)
	})
	userGroup.GET("/get_recommendation", func(c *gin.Context) {
		coreApi.UserApi.GetRecommendation(c, c.Request)
	})
}

func connectToDatabase(cnf config.ServiceConfiguration) *sqlx.DB {
	psqlInfo := cnf.PostgresDSN.String()
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func setUpApi(db *sqlx.DB, client *http.Client, cnf config.ServiceConfiguration, jwtManager *auth.Manager) *api.Api {
	postgresStorage := storage.NewPostgresStorage(db)
	telethonService := services.NewTelethonService(client)

	authService := auth.NewAuthService(postgresStorage, telethonService, jwtManager)
	userService := user.NewUserService(postgresStorage, telethonService)
	feedService := feed.NewFeedService(postgresStorage, telethonService)

	userApi := api.NewUserAPI(userService)
	authApi := api.NewAuthAPI(authService)
	feedApi := api.NewFeedAPI(feedService)
	coreApi := api.NewAPI(authApi, userApi, feedApi, cnf)
	return coreApi
}

func main() {

	cnf := config.NewConfig()
	jwtManager := auth.NewManger(cnf)
	db := connectToDatabase(cnf)
	client := &http.Client{
		Timeout: time.Minute,
	}
	coreApi := setUpApi(db, client, cnf, jwtManager)
	r := gin.Default()
	r.Use(CORSMiddleware())
	setRoutes(r, coreApi)
	r.Run(":8082")

	defer db.Close()
}
