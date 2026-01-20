package startgg

type RawPhaseGroupsData struct {
	Data   DataEvent `json:"data"`
	Errors []Errors  `json:"errors"`
}

type DataEvent struct {
	Event Event `json:"event"`
}

type Event struct {
	Id          int64          `json:"id"`
	Name        string         `json:"name"`
	State       StateEvent     `json:"state"`
	PhaseGroups []PhaseGroupId `json:"phaseGroups"`
}

type PhaseGroupId struct {
	Id int64 `json:"id"`
}

func (c *Client) GetListGroups(slug string) ([]PhaseGroupId, error) {
	var variables = map[string]any{
		"slug": slug,
	}

	results, err := GetData[RawPhaseGroupsData](c, getListPhaseGroups, variables)
	if err != nil {
		return nil, err
	}

	return results.Data.Event.PhaseGroups, nil
}
