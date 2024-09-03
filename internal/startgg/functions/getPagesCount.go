package functions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

type RawPagesDataCount struct {
	Data   DataPhaseGroup   `json:"data"`
	Errors []startgg.Errors `json:"errors"`
}

func GetPagesCount(phaseGroupID int64) (int, error) {
	if !startgg.Token() {
		return 0, errors.New("token verification - authentication token not set")
	}

	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
	}

	query, err := json.Marshal(startgg.PrepareQuery(startgg.GetPagesCount, variables))
	if err != nil {
		return 0, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := startgg.RunQuery(query)
	if err != nil {
		return 0, err
	}

	results := &RawPagesDataCount{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return 0, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results.Data.PhaseGroup.Sets.PageInfo.Total, nil
}
