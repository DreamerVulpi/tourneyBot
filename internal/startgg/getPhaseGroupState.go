package startgg

import (
	"encoding/json"
	"fmt"
)

type RawPhaseGroupStateData struct {
	Data   DataPhaseGroupState `json:"data"`
	Errors []Errors            `json:"errors"`
}

type DataPhaseGroupState struct {
	PhaseGroup PGState `json:"phaseGroup"`
}

type PGState struct {
	Id    int64 `json:"id"`
	State int   `json:"state"`
}

func (c *Client) GetPhaseGroupState(phaseGroupID int64) (int, error) {
	// if !startgg.Token() {
	// 	return 0, errors.New("token Verification - Authentication Token Not Set")
	// }

	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
	}

	query, err := json.Marshal(PrepareQuery(getPhaseGroupState, variables))
	if err != nil {
		return 0, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := c.RunQuery(query)
	if err != nil {
		return 0, err
	}

	results := &RawPhaseGroupStateData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return 0, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results.Data.PhaseGroup.State, nil
}