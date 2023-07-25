package models

import "encoding/json"

type Channel struct {
	Name string `json:"name" db:"channel_name"`
	Tag  string `json:"tag" db:"channel_id"`
}

func (codeRequest GetChannelsRequest) Marshall() ([]byte, error) {
	jsonBytes, err := json.Marshal(codeRequest)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

type GetChannelsRequest struct {
	UserID  uint32 `json:"id"`
	Command string `json:"command"`
}

type GetPublicUserChannelsResult struct {
	Command string    `json:"command"`
	Details []Channel `json:"details"`
	ID      int       `json:"id"`
}

type FormFeedRequest struct {
	UserID             uint32   `json:"id"`
	Command            string   `json:"command"`
	MessagesLimitCount int32    `json:"mes_count"`
	ChannelsIDs        []string `json:"chat_ids"`
}

type GetRecommendationRequest struct {
	UserID  uint32 `json:"id"`
	Command string `json:"command"`
}

type GetRecommendationRequestResult struct {
	Command string   `json:"command"`
	Details []string `json:"details"`
	UserID  uint32   `json:"id"`
}

//type ChannelsList struct {
//}
