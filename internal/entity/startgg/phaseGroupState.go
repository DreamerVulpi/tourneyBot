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
