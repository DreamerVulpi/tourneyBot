package startgg

type RawPagesDataCount struct {
	Data   DataPhaseGroupSets `json:"data"`
	Errors []Errors           `json:"errors"`
}
