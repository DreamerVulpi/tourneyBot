package entity

import "time"

type Participant struct {
	GamerTag               string    `json:"gamerTag"`
	MessengerPlatform      string    `json:"messengerPlatform"`
	MessengerPlatformId    string    `json:"messengerPlatformId"`
	MessengerPlatformLogin string    `json:"messengerPlatformLogin"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type ParticipantAddRequest struct {
	GamerTag               string    `json:"gamerTag"`
	MessengerPlatform      string    `json:"messengerPlatform"`
	MessengerPlatformId    string    `json:"messengerPlatformId"`
	MessengerPlatformLogin string    `json:"messengerPlatformLogin"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type ParticipantEditRequest struct {
	GamerTag                string    `json:"gamerTag"`
	MessenagerPlatform      string    `json:"messenagerPlatform"`
	MessenagerPlatformId    string    `json:"messenagerPlatformId"`
	MessenagerPlatformLogin string    `json:"messenagerPlatformLogin"`
	UpdatedAt               time.Time `json:"updatedAt"`
}

type ParticipantDeleteRequest struct{}

type ParticipantAddResponse struct {
	GamerTag          string `json:"gamerTag"`
	MessengerPlatform string `json:"messengerPlatform"`
}

type ParticipantEditResponse struct{}
type ParticipantDeleteResponse struct{}
type ParticipantGetResponse struct {
	GamerTag               string    `json:"gamerTag"`
	MessengerPlatform      string    `json:"messengerPlatform"`
	MessengerPlatformId    string    `json:"messengerPlatformId"`
	MessengerPlatformLogin string    `json:"messengerPlatformLogin"`
	UpdatedAt              time.Time `json:"updatedAt"`
}
