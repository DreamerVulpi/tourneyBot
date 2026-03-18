package usecase

import (
	"context"

	"github.com/dreamervulpi/tourneyBot/internal/sender"
	"github.com/dreamervulpi/tourneyBot/startgg"
)

type StartggMatchAdapter struct {
	Source         []startgg.Nodes
	TournamentName string
	TournamentSlug string
	PhaseGroupId   int64
	IsFinals       bool
	DebugMode      bool
	TestUser       sender.Participant
	Contacts       map[string]sender.Participant
}

func (a StartggMatchAdapter) GetSetsData(ctx context.Context, ns sender.NotificationSender) ([]sender.SetData, error) {
	return nil, nil
}
