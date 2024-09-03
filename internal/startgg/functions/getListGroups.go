package functions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type RawPhaseGroupsData struct {
	Data   DataEvent        `json:"data"`
	Errors []startgg.Errors `json:"errors"`
}

type DataEvent struct {
	Event Event `json:"event"`
}

type Event struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	PhaseGroups []PGS  `json:"phaseGroups"`
}

type PGS struct {
	Id int64 `json:"id"`
}

func GetListGroups(slug string) ([]PGS, error) {
	if !startgg.Token() {
		return []PGS{}, errors.New("token verification - authentication token not set")
	}

	var variables = map[string]any{
		"slug": slug,
	}

	query, err := json.Marshal(startgg.PrepareQuery(startgg.GetListPhaseGroups, variables))
	if err != nil {
		return []PGS{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := startgg.RunQuery(query)
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
