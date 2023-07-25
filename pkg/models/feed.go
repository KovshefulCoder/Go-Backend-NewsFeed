package models

type FormFeedResponse struct {
	Command string      `json:"command"`
	Details []GroupNews `json:"details"`
	ID      uint32      `json:"id"`
}

type GroupNews struct {
	Group     []News `json:"group"`
	GroupType uint8  `json:"type"`
}

type News struct {
	Date          string `json:"data"`
	SourceChannel string `json:"source_channel"`
	Text          string `json:"text"`
}
