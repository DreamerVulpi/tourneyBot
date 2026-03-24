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
