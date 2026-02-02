package sender

import (
	"context"
)

type Participant struct {
	DiscordID    string
	DiscordLogin string
	GameNickname string
	GameID       string
	Locales      []string
}

type MatchInfo struct {
}

type Sender interface {
	SendInvite(ctx context.Context, p Participant, match MatchInfo) error
}

type TournamentPlatform interface {
	GetActiveMatches(slug string) ([]MatchInfo, error)
}
