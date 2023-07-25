package feed

import (
	"NewsFeed/pkg/models"
	"NewsFeed/pkg/services"
	"log"
)

type Storage interface {
	GetUserFeedChannels(userID uint32) ([]models.Channel, error)
	AddFeedChannel(userID uint32, channelTag string, channelName string) error
	RemoveFeedChannel(userID uint32, channelTag string) error
}

type Service struct {
	storage Storage
	th      *services.TelethonService
}

func NewFeedService(storage Storage, th *services.TelethonService) *Service {
	return &Service{storage: storage, th: th}
}

func (s *Service) GetFeedChannels(userID uint32) ([]models.Channel, error) {
	channels, err := s.storage.GetUserFeedChannels(userID)
	if err != nil {
		log.Printf("GetFeedChannels: can`t retrieve user`s channels from feed: %v", err)
		return nil, err
	}
	return channels, nil
}

func (s *Service) AddFeedChannel(userID uint32, channelTag string, channelName string) error {
	err := s.storage.AddFeedChannel(userID, channelTag, channelName)
	if err != nil {
		log.Printf("AddFeedChannel: can`t add channel to feed: %v", err)
		return err
	}
	return nil
}

func (s *Service) RemoveFeedChannel(userID uint32, channelTag string) error {
	err := s.storage.RemoveFeedChannel(userID, channelTag)
	if err != nil {
		log.Printf("RemoveFeedChannel: can`t remove channel from feed: %v", err)
		return err
	}
	return nil
}
