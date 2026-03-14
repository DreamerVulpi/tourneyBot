package entity

import "time"

type SentSet struct {
	SetId              int64     `json:"setId"`
	TournamentPlatform string    `json:"tournamentPlatform"`
	MessengerPlatform  string    `json:"messengerPlatform"`
	TournamentSlug     string    `json:"tournamentSlug"`
	SentAt             time.Time `json:"sentat"`
}

type SentSetCheckRequest struct {
	SetId int64 `json:"setId"`
}

type SentSetAddRequest struct {
	SetId              int64     `json:"setId"`
	TournamentPlatform string    `json:"sourcePlatform"`
	MessengerPlatform  string    `json:"messengerPlatform"`
	TournamentSlug     string    `json:"tournamentSlug"`
	SentAt             time.Time `json:"sentat"`
}

type SentSetEditRequest struct {
	SetId  int64     `json:"setId"`
	SentAt time.Time `json:"sentat"`
}

type SentSetDeleteRequest struct {
	SetId int64 `json:"setId"`
}

type SentSetGetRequest struct {
	SetId int64 `json:"setId"`
}

type SentSetAddResponse struct {
	SetId int64 `json:"setId"`
}

type SentSetCheckResponse struct {
	State bool `json:"state"`
}

type SentSetEditResponse struct{}

type SentSetDeleteResponse struct{}

type SentSetGetResponse struct {
	SetId              int64     `json:"setId"`
	TournamentPlatform string    `json:"tournamentPlatform"`
	MessengerPlatform  string    `json:"messengerPlatform"`
	TournamentSlug     string    `json:"tournamentSlug"`
	SentAt             time.Time `json:"sentat"`
}
