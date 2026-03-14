package sender

import (
	"context"
)

type Participant struct {
	MessenagerID    string
	MessenagerLogin string
	GameNickname    string
	GameID          string
	Locales         []string
}

type SetData struct {
	TournamentName string
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
	SendNotification(ctx context.Context, targetID string, data SetData) error
	GetPlatformMessenagerName() string
}
