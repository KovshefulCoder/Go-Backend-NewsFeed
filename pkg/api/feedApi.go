package api

import (
	"NewsFeed/pkg/services/feed"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FeedAPI struct {
	FeedService *feed.Service
}

func NewFeedAPI(FeedService *feed.Service) *FeedAPI {
	return &FeedAPI{FeedService: FeedService}
}

func (api *FeedAPI) GetFeedChannels(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodGet {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	feedChannels, err := api.FeedService.GetFeedChannels(userID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, feedChannels)
}

func (api *FeedAPI) AddFeedChannel(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var addFeedChannelBody ChangeFeedChannelBody
	err = json.NewDecoder(r.Body).Decode(&addFeedChannelBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	err = api.FeedService.AddFeedChannel(userID, addFeedChannelBody.Tag, addFeedChannelBody.Name)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

func (api *FeedAPI) RemoveFeedChannel(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodDelete {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var removeFeedChannelBody ChangeFeedChannelBody
	err = json.NewDecoder(r.Body).Decode(&removeFeedChannelBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	err = api.FeedService.RemoveFeedChannel(userID, removeFeedChannelBody.Tag)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

type ChangeFeedChannelBody struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}
