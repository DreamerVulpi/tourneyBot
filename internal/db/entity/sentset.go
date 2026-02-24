package entity

import "time"

type SentSet struct {
	SetId              int64     `json:"setId"`
	TournamentPlatform string    `json:"tournamentPlatform"`
	MessengerPlatform  string    `json:"messengerPlatform"`
	TournamentSlug     string    `json:"tournamentSlug"`
	SentAt             time.Time `json:"sentat"`
}

type SentSetAddRequest struct {
	SetId              int64     `json:"setId"`
	TournamentPlatform string    `json:"sourcePlatform"`
	MessengerPlatform  string    `json:"messengerPlatform"`
	TournamentSlug     string    `json:"tournamentSlug"`
	SentAt             time.Time `json:"sentat"`
}

type SentSetEditRequest struct {
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

type SentSetEditResponse struct{}

type SentSetDeleteResponse struct{}

type SentSetGetResponse struct {
	SetId              int64     `json:"setId"`
	TournamentPlatform string    `json:"tournamentPlatform"`
	MessengerPlatform  string    `json:"messengerPlatform"`
	TournamentSlug     string    `json:"tournamentSlug"`
	SentAt             time.Time `json:"sentat"`
}
