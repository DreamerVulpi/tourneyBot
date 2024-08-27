package startgg

import (
	"encoding/json"
	"errors"
	"fmt"
)

type RawPhaseGroupStateData struct {
	Data   DataPhaseGroupState `json:"data"`
	Errors []Errors            `json:"errors"`
}

type DataPhaseGroupState struct {
	PhaseGroup PGState `json:"phaseGroup"`
}

type PGState struct {
	Id    int64 `json:"id"`
	State int   `json:"state"`
}

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
	Id    int64   `json:"id"`
	State int     `json:"state"`
	Slots []Slots `json:"slots"`
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

func GetPhaseGroupState(phaseGroupID int64) (*RawPhaseGroupStateData, error) {
	if !token() {
		return &RawPhaseGroupStateData{}, errors.New("Token Verification - Authentication Token Not Set")
	}

	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
	}

	query, err := json.Marshal(prepareQuery(getPhaseGroupState, variables))
	if err != nil {
		return &RawPhaseGroupStateData{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := runQuery(query)
	if err != nil {
		return &RawPhaseGroupStateData{}, err
	}

	results := &RawPhaseGroupStateData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results, nil
}

func GetPhaseGroupSets(phaseGroupID int64, page int, perPage int) (*RawPhaseGroupData, error) {
	if !token() {
		return &RawPhaseGroupData{}, errors.New("token verification - authentication token not set")
	}

	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
		"page":         page,
		"perPage":      perPage,
	}

	query, err := json.Marshal(prepareQuery(getPhaseGroupSets, variables))
	if err != nil {
		return &RawPhaseGroupData{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := runQuery(query)
	if err != nil {
		return &RawPhaseGroupData{}, err
	}

	results := &RawPhaseGroupData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results, nil
}
