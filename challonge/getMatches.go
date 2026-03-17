package challonge

import (
	"context"
	"fmt"
)

// State of the Matches
type State string

const (
	Complete State = "complete"
	Open     State = "open"
	Pending  State = "pending"
)

type ApidogMatchesResponse struct {
	Data []ModelMatches `json:"data"`
}

type ModelMatches struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Attributes    MatchOutput        `json:"attributes"`
	Relationships MatchRelationships `json:"relationships"`
}

type MatchOutput struct {
	Identifier          string             `json:"identifier"`
	PointsByParticipant []PointRecord      `json:"pointsByParticipant"`
	Round               int64              `json:"round"`
	ScoreInSets         [][]int            `json:"scoreInSets"`
	Scores              string             `json:"scores"`
	State               State              `json:"state"`
	SuggestedPlayOrder  int64              `json:"suggestedPlayOrder"`
	Timestamps          Timestamps         `json:"timestamps"`
	WinnerID            int64              `json:"winners"`
	Relationships       MatchRelationships `json:"-"`
}

type Timestamps struct {
	StartedAt  string `json:"startedAt"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	UnderwayAt string `json:"underwayAt"`
}

type PointRecord struct {
	ParticipantID interface{} `json:"participantId"`
	Scores        []int       `json:"scores"`
}

type MatchRelationships struct {
	Player1 *PlayerRelation `json:"player1"`
	Player2 *PlayerRelation `json:"player2"`
}

type PlayerRelation struct {
	Data *PlayerData `json:"data"`
}

type PlayerData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (c *Client) GetMatches(ctx context.Context, tourneySlug string) ([]MatchOutput, error) {
	slug := ExtractSlug(tourneySlug)

	results, err := GetData[ApidogMatchesResponse](c, ctx, getMatches, slug)
	if err != nil {
		return nil, fmt.Errorf("getMatches | Error getting data: %w", err)
	}

	var matches []MatchOutput
	for _, m := range results.Data {
		attr := m.Attributes
		attr.Relationships = m.Relationships
		matches = append(matches, attr)
	}
	return matches, nil
}

func (m MatchOutput) String() string {
	return fmt.Sprintf("[%s] %s vs %s | Winner: %v | State: %s\n",
		m.Identifier,
		m.Relationships.Player1.Data.ID,
		m.Relationships.Player2.Data.ID,
		m.WinnerID,
		m.State,
	)
}
