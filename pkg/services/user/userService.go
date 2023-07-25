package user

import (
	"NewsFeed/pkg/models"
	"NewsFeed/pkg/services"
	"log"
	"sort"
	"time"
)

type Storage interface {
	GetUserFeedChannels(userID uint32) ([]models.Channel, error)
}

type Service struct {
	storage Storage
	th      *services.TelethonService
}

func NewUserService(storage Storage, th *services.TelethonService) *Service {
	return &Service{storage: storage, th: th}
}

func (s *Service) GetAllPublicSubsChannels(userID uint32) ([]models.Channel, error) {
	getPublicChannelsBody := models.GetChannelsRequest{UserID: userID, Command: "get_chats"}
	channels, err := s.th.GetPublicUserChannels(getPublicChannelsBody)
	if err != nil {
		log.Printf("GetAllPublicSubsChannels: can`t retrieve user`s public channels subs: %v", err)
		return nil, err
	}
	return channels, nil
}

func (s *Service) FormFeed(userID uint32) ([]models.GroupNews, error) {
	channels, err := s.storage.GetUserFeedChannels(userID)
	if err != nil {
		log.Printf("FormFeed: can`t retrieve user`s feed channels: %v", err)
		return nil, err
	}
	channelsIDs := make([]string, 0)
	for _, channel := range channels {
		channelsIDs = append(channelsIDs, channel.Tag)
	}
	formFeedBody := models.FormFeedRequest{UserID: userID, Command: "get_messages", MessagesLimitCount: 100, ChannelsIDs: channelsIDs}
	feed, err := s.th.FormFeed(formFeedBody)
	if err != nil {
		log.Printf("FormFeed: can`t form feed: %v", err)
		return nil, err
	}
	sort.Slice(feed.Details, func(i, j int) bool {
		// Extract the earliest date from each GroupNews
		iDate := extractEarliestDate(feed.Details[i])
		jDate := extractEarliestDate(feed.Details[j])

		// Parse the dates
		iTime, _ := time.Parse(time.RFC1123, iDate)
		jTime, _ := time.Parse(time.RFC1123, jDate)

		// Compare the dates
		return iTime.After(jTime)
	})
	return feed.Details, nil
}

func (s *Service) GetRecommendation(userID uint32) ([]string, error) {
	channels, err := s.storage.GetUserFeedChannels(userID)
	if err != nil {
		log.Printf("GetRecommendation: can`t retrieve user`s feed channels: %v", err)
		return nil, err
	}
	channelsIDs := make([]string, len(channels))
	for _, channel := range channels {
		channelsIDs = append(channelsIDs, channel.Tag)
	}
	getRecommendationBody := models.GetRecommendationRequest{UserID: userID, Command: "get_recomendations"}
	recommendation, err := s.th.GetRecommendation(getRecommendationBody)
	if err != nil {
		log.Printf("GetRecommendation: can`t get recommendation: %v", err)
		return nil, err
	}
	return recommendation, nil
}

func extractEarliestDate(group models.GroupNews) string {
	earliestDateString := ""
	if group.GroupType == 0 {
		earliestDateString = group.Group[0].Date
	} else {
		sort.Slice(group.Group, func(i, j int) bool {
			// Extract the earliest date from each GroupNews
			iDate := group.Group[i].Date
			jDate := group.Group[j].Date

			// Parse the dates
			iTime, _ := time.Parse(time.RFC1123, iDate)
			jTime, _ := time.Parse(time.RFC1123, jDate)

			// Compare the dates
			return iTime.Before(jTime)
		})
		earliestDateString = group.Group[0].Date
	}
	return earliestDateString
}
