package sender

import (
	"context"
)

type Participant struct {
	MessenagerID    string
	MessenagerLogin string
	MessenagerName  string
	GameNickname    string
	GameID          string
	Locales         []string
}

type SetData struct {
	TournamentName string
	// TODO: Change to string
	SetID          int64
	StreamName     string
	StreamSourse   string
	RoundNum       int
	PhaseGroupId   int64
	Recipient      Participant
	Opponent       Participant
	FullInviteLink string
	IsFinals       bool
	IsTest         bool
}

type NotificationSender interface {
	FindContactOfParticipant(ctx context.Context, participant Participant) (Participant, error)
	SendNotification(ctx context.Context, targetID string, data SetData) error
	GetPlatformMessenagerName() string
}

type NotificationData interface {
	// CheckContactOfParticipant(participant Participant) (Participant, error)
	GetSetsData(ctx context.Context, ns NotificationSender) ([]SetData, error)
	GetPlatformTournamentName() string
}
