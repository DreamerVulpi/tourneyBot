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

func (c *Client) GetTournament(tourneySlug string) (Tournament, error) {
	var variables = map[string]any{
		"tourneySlug": tourneySlug,
	}

	results, err := GetData[RawTournamentData](c, getTournament, variables)
	if err != nil {
		return Tournament{}, err
	}

	return results.Data.Tournament, nil
}
