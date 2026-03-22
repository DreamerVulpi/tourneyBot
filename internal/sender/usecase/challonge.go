package usecase

import (
	"context"

	"github.com/dreamervulpi/tourneyBot/challonge"
	"github.com/dreamervulpi/tourneyBot/internal/sender"
)

type ChallongeMatchAdapter struct {
	Client         *challonge.Client
	TournamentSlug string
	Finals         StartggFinalConfig
	DebugMode      bool
	TestUser       sender.Participant
}

func (c ChallongeMatchAdapter) GetPlatformTournamentName() string {
	return "challonge"
}

func (c ChallongeMatchAdapter) GetTournamentSlug() string {
	return c.TournamentSlug
}

func (c ChallongeMatchAdapter) GetSetsData(ctx context.Context) ([]sender.SetData, error) {
	// TODO: Complete method
	return nil, nil
}
