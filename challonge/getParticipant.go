package challonge

import (
	"context"
	"fmt"
)

type ApidogParticipantResponse struct {
	Data ParticipantModel `json:"data"`
}

type ParticipantModel struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes ParticipantOutput `json:"attributes"`
}

// ParticipantOutput
type ParticipantOutput struct {
	ID           string     `json:"-"`
	FinalRank    int64      `json:"finalRank"`
	GroupID      int64      `json:"groupId"`
	Misc         string     `json:"misc"`
	Name         string     `json:"name"`
	Seed         int64      `json:"seed"`
	States       States     `json:"states"`
	Timestamps   Timestamps `json:"timestamps"`
	TournamentID int64      `json:"tournamentId"`
	Username     string     `json:"username"`
	LinkCheckIn  string     `json:"-"`
}

type States struct {
	Active bool `json:"active"`
}

func (c *Client) GetParticipant(ctx context.Context, tourneySlug string, participantId string) (ParticipantOutput, error) {
	slug := ExtractSlug(tourneySlug)

	result, err := GetData[ApidogParticipantResponse](c, ctx, getParticipant, slug, participantId)
	if err != nil {
		return ParticipantOutput{}, fmt.Errorf("getParticipant | Error getting data: %w", err)
	}

	return ParticipantOutput{
		ID:           result.Data.ID,
		FinalRank:    result.Data.Attributes.FinalRank,
		GroupID:      result.Data.Attributes.GroupID,
		Misc:         result.Data.Attributes.Misc,
		Name:         result.Data.Attributes.Name,
		Seed:         result.Data.Attributes.Seed,
		States:       result.Data.Attributes.States,
		Timestamps:   result.Data.Attributes.Timestamps,
		TournamentID: result.Data.Attributes.TournamentID,
		Username:     result.Data.Attributes.Username,
		LinkCheckIn:  fmt.Sprintf(getCheckIn, slug, result.Data.ID),
	}, nil
}
