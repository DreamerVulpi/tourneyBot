package challonge

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/dreamervulpi/tourneyBot/challonge"
	"github.com/stretchr/testify/assert"
)

func TestGetDataChallonge(t *testing.T) {
	rawData, err := os.ReadFile("rawDataMatches.json")
	if err != nil {
		t.Fatalf("Failed to read rawDataMatches.json: %v", err)
	}

	expectedFirstMatch := challonge.ModelMatches{
		ID:   "186702200",
		Type: "match",
		Attributes: challonge.MatchOutput{
			State:              challonge.Complete,
			Round:              1,
			Identifier:         "A",
			Scores:             "0 - 2",
			SuggestedPlayOrder: 1,
			ScoreInSets:        [][]int{{0, 2}},
			PointsByParticipant: []challonge.PointRecord{
				{ParticipantID: float64(112579135), Scores: []int{0}},
				{ParticipantID: float64(112579132), Scores: []int{2}},
			},
			Timestamps: challonge.Timestamps{
				StartedAt:  "2020-01-11T14:00:50.421Z",
				CreatedAt:  "2020-01-11T14:00:50.273Z",
				UpdatedAt:  "2020-01-11T14:31:49.537Z",
				UnderwayAt: "",
			},
			WinnerID: 112579132,
		},
		Relationships: challonge.MatchRelationships{
			Player1: &challonge.PlayerRelation{
				Data: &challonge.PlayerData{
					ID:   "112579135",
					Type: "participant",
				},
			},
			Player2: &challonge.PlayerRelation{
				Data: &challonge.PlayerData{
					ID:   "112579132",
					Type: "participant",
				},
			},
		},
	}

	var actualResponse challonge.ApidogMatchesResponse
	err = json.Unmarshal(rawData, &actualResponse)
	if err != nil {
		assert.NoError(t, fmt.Errorf("JSON Unmarshal failed - %w", err))
	}

	assert.Equal(t, 8, len(actualResponse.Data))

	actualFirstMatch := actualResponse.Data[0]

	assert.Equal(t, expectedFirstMatch.ID, actualFirstMatch.ID)
	assert.Equal(t, expectedFirstMatch.Attributes.Identifier, actualFirstMatch.Attributes.Identifier)
	assert.Equal(t, expectedFirstMatch.Attributes.State, actualFirstMatch.Attributes.State)
	assert.Equal(t, expectedFirstMatch.Attributes.WinnerID, actualFirstMatch.Attributes.WinnerID)

	assert.NotNil(t, actualFirstMatch.Relationships.Player1)
	assert.Equal(t, "112579135", actualFirstMatch.Relationships.Player1.Data.ID)
}

func TestExtractSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://challonge.com/tournament123", "tournament123"},
		{"https://challonge.com/tournament123/", "tournament123"},
		{"my_tourney", "my_tourney"},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expected, challonge.ExtractSlug(tc.input))
	}
}
