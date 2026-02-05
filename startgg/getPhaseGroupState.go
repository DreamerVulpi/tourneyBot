package startgg

type RawPhaseGroupStateData struct {
	Data   DataPhaseGroupState `json:"data"`
	Errors []Errors            `json:"errors"`
}

type DataPhaseGroupState struct {
	PhaseGroup PhaseGroupState `json:"phaseGroup"`
}

type PhaseGroupState struct {
	Id    int64 `json:"id"`
	State int   `json:"state"`
}

func (c *Client) GetPhaseGroupState(phaseGroupID int64) (State, error) {
	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
	}

	results, err := GetData[RawPhaseGroupStateData](c, getPhaseGroupState, variables)
	if err != nil {
		return 0, err
	}

	return State(results.Data.PhaseGroup.State), nil
}
