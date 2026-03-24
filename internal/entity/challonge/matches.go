package challonge

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
	ID                  string             `json:"-"`
	Identifier          string             `json:"identifier"`
	PointsByParticipant []PointRecord      `json:"pointsByParticipant"`
	Round               int                `json:"round"`
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
	ParticipantID string `json:"participantId"`
	Scores        []int  `json:"scores"`
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
