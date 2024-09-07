package startgg

import (
	"encoding/json"
	"fmt"
)

type RawPagesDataCount struct {
	Data   DataPhaseGroup `json:"data"`
	Errors []Errors       `json:"errors"`
}

func (c *Client) GetPagesCount(phaseGroupID int64) (int, error) {
	// if !startgg.Token() {
	// 	return 0, errors.New("token verification - authentication token not set")
	// }

	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
	}

	query, err := json.Marshal(PrepareQuery(getPagesCount, variables))
	if err != nil {
		return 0, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := c.RunQuery(query)
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
