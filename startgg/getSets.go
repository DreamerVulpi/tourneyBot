package startgg

type RawPhaseGroupData struct {
	Data   DataPhaseGroup `json:"data"`
	Errors []Errors       `json:"errors"`
}

type DataPhaseGroup struct {
	PhaseGroup PhaseGroupSets `json:"phaseGroup"`
}

// Group(Phase)
type PhaseGroupSets struct {
	Id   int64 `json:"id"`
	Sets Sets  `json:"sets"`
}

type Sets struct {
	PageInfo PageInfo `json:"pageInfo"`
	Nodes    []Nodes  `json:"nodes"`
}

// Sets counts in Group(Phase)
type PageInfo struct {
	Total int `json:"total"`
}

// Information about Set
type Nodes struct {
	Id            int64    `json:"id"`
	State         State    `json:"state"`
	FullRoundText string   `json:"fullRoundText"`
	Round         int      `json:"round"`
	Stream        Streamer `json:"stream"`
	Slots         []Slots  `json:"slots"`
}

type Streamer struct {
	StreamName   string `json:"streamName"`
	StreamSource string `json:"streamSource"`
}

// Slots in set
type Slots struct {
	Entrant Entrant `json:"entrant"`
}

// Player in tournament
type Entrant struct {
	Id           int64          `json:"id"`
	Participants []Participants `json:"participants"`
}

// Player data
type Participants struct {
	GamerTag          string            `json:"gamerTag"`
	ConnectedAccounts ConnectedAccounts `json:"connectedAccounts"`
	User              User              `json:"user"`
}

// Linked accounts
type ConnectedAccounts struct {
	Tekken Tekken8 `json:"tekken"`
	SF6    SF6     `json:"capcom"`
}

type Tekken8 struct {
	TekkenID string `json:"value"`
}

type SF6 struct {
	GameID string `json:"value"`
}

// User authorizations (Discord and etc.)
type User struct {
	Authorizations []Authorizations `json:"authorizations"`
}

type Authorizations struct {
	Discord string `json:"externalUsername"`
}

func (c *Client) GetSets(phaseGroupID int64, page int, perPage int) ([]Nodes, error) {
	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
		"page":         page,
		"perPage":      perPage,
	}

	// GetPhaseGroupSets || testGetPhaseGroupSets
	results, err := GetData[RawPhaseGroupData](c, testGetPhaseGroupSets, variables)
	if err != nil {
		return nil, err
	}

	return results.Data.PhaseGroup.Sets.Nodes, nil
}
