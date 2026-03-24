package startgg

type RawTournamentData struct {
	Data   DataTournament `json:"data"`
	Errors []Errors       `json:"errors"`
}

type DataTournament struct {
	Tournament Tournament `json:"tournament"`
}

type Tournament struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	State State  `json:"state"`
}
