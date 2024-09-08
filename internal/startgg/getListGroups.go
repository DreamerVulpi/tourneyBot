package startgg

import (
	"encoding/json"
	"fmt"
)

type RawPhaseGroupsData struct {
	Data   DataEvent `json:"data"`
	Errors []Errors  `json:"errors"`
}

type DataEvent struct {
	Event Event `json:"event"`
}

type Event struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	State       StateEvent `json:"state"`
	PhaseGroups []PGS      `json:"phaseGroups"`
}

type PGS struct {
	Id int64 `json:"id"`
}

func (c *Client) GetListGroups(slug string) ([]PGS, error) {
	var variables = map[string]any{
		"slug": slug,
	}

	query, err := json.Marshal(PrepareQuery(getListPhaseGroups, variables))
	if err != nil {
		return []PGS{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := c.RunQuery(query)
	if err != nil {
		return []PGS{}, err
	}

	results := &RawPhaseGroupsData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results.Data.Event.PhaseGroups, nil
}
