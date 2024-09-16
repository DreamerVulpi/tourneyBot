package startgg

import (
	"encoding/json"
	"fmt"
)

type RawTournamentData struct {
	Data   DataTournament `json:"data"`
	Errors []Errors       `json:"errors"`
}

type DataTournament struct {
	Tournament Tournament `json:"tournament"`
}

type Tournament struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	State State  `json:"state"`
}

func (c *Client) GetTournament(tourneySlug string) (Tournament, error) {
	var variables = map[string]any{
		"tourneySlug": tourneySlug,
	}

	query, err := json.Marshal(PrepareQuery(getTournament, variables))
	if err != nil {
		return Tournament{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := c.RunQuery(query)
	if err != nil {
		return Tournament{}, err
	}

	results := &RawTournamentData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return Tournament{}, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results.Data.Tournament, nil
}
