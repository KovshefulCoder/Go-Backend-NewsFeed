package api

import (
	"NewsFeed/pkg/services/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UserAPI struct {
	UserService *user.Service
}

func NewUserAPI(UserService *user.Service) *UserAPI {
	return &UserAPI{UserService: UserService}
}

func (api *UserAPI) GetAllPublicSubsChannels(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodGet {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("UserAPI GetAllPublicSubsChannels: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	channels, err := api.UserService.GetAllPublicSubsChannels(userID)
	if err != nil {
		log.Println("UserAPI GetAllPublicSubsChannels: can`t get user`s public channels from subscriptions")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, channels)
}

func (api *UserAPI) FormFeed(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodGet {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("UserAPI FormFeed: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	feed, err := api.UserService.FormFeed(userID)
	if err != nil {
		log.Println("UserAPI FormFeed: can`t form feed")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, feed)
}

func (api *UserAPI) GetRecommendation(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodGet {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("UserAPI GetRecommendation: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	recommendation, err := api.UserService.GetRecommendation(userID)
	if err != nil {
		log.Println("UserAPI GetRecommendation: can`t get recommendation")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, recommendation)
}
