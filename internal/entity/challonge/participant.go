package challonge

type ApidogParticipantResponse struct {
	Data ParticipantModel `json:"data"`
}

type ParticipantModel struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes ParticipantOutput `json:"attributes"`
}

// ParticipantOutput
type ParticipantOutput struct {
	ID           string     `json:"-"`
	FinalRank    int64      `json:"finalRank"`
	GroupID      int64      `json:"groupId"`
	Misc         string     `json:"misc"`
	Name         string     `json:"name"`
	Seed         int64      `json:"seed"`
	States       States     `json:"states"`
	Timestamps   Timestamps `json:"timestamps"`
	TournamentID int64      `json:"tournamentId"`
	Username     string     `json:"username"`
	LinkCheckIn  string     `json:"-"`
}

type States struct {
	Active bool `json:"active"`
}
