package startgg

import (
	"encoding/json"
	"fmt"
)

type RawPhaseGroupData struct {
	Data   DataPhaseGroup `json:"data"`
	Errors []Errors       `json:"errors"`
}

type DataPhaseGroup struct {
	PhaseGroup PhaseGroup `json:"phaseGroup"`
}

// Group(Phase)
type PhaseGroup struct {
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
	Id     int64    `json:"id"`
	State  int      `json:"state"`
	Stream Streamer `json:"stream"`
	Slots  []Slots  `json:"slots"`
}

type Streamer struct {
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

type Participants struct {
	GamerTag          string            `json:"gamerTag"`
	ConnectedAccounts ConnectedAccounts `json:"connectedAccounts"`
	User              User              `json:"user"`
}

type ConnectedAccounts struct {
	Tekken Tekken8 `json:"tekken"`
}

type Tekken8 struct {
	TekkenID string `json:"value"`
}

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

	query, err := json.Marshal(PrepareQuery(getPhaseGroupSets, variables))
	if err != nil {
		return []Nodes{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := c.RunQuery(query)
	if err != nil {
		return []Nodes{}, err
	}

	results := &RawPhaseGroupData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results.Data.PhaseGroup.Sets.Nodes, nil
}
