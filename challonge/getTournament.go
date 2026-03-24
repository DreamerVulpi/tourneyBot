package challonge

import (
	"context"
	"fmt"
)

type ApidogTournamentResponse struct {
	Data *TournamentModel `json:"data"`
}

type TournamentModel struct {
	ID         *string     `json:"id"`
	Type       *string     `json:"type"`
	Attributes *Tournament `json:"attributes"`
}

type Tournament struct {
	Description              *string                   `json:"description"`
	DoubleEliminationOptions *DoubleEliminationOptions `json:"double_elimination_options"`
	FreeForAllOptions        *FreeForAllOptions        `json:"free_for_all_options"`
	GameName                 *string                   `json:"game_name"`
	// When true, Challonge's two-stage format is used (group stage + final stage). Additional
	// tournament state transitions must be performed (start_group_stage, finalize_group_stage).
	GroupStageEnabled *bool              `json:"group_stage_enabled"`
	GroupStageOptions *GroupStageOptions `json:"group_stage_options"`
	MatchOptions      *MatchOptions      `json:"match_options"`
	Name              string             `json:"name"`
	Notifications     *Notifications     `json:"notifications"`
	// When true, this tournament will be opted out for search engine bots and will be hidden
	// from indexes on Challonge
	Private             *bool                `json:"private"`
	RegistrationOptions *RegistrationOptions `json:"registration_options"`
	RoundRobinOptions   *RoundRobinOptions   `json:"round_robin_options"`
	SeedingOptions      *SeedingOptions      `json:"seeding_options"`
	// The scheduled start time. Note that Challonge will not automatically start the tournament.
	StartsAt       *string         `json:"starts_at"`
	StationOptions *StationOptions `json:"station_options"`
	SwissOptions   *SwissOptions   `json:"swiss_options"`
	TournamentType TournamentType  `json:"tournament_type"`
	URL            *string         `json:"url"`
}

type DoubleEliminationOptions struct {
	// When blank, the losers bracket winner has to beat the winners bracket winner twice.
	GrandFinalsModifier *GrandFinalsModifierEnum `json:"grand_finals_modifier"`
	// Double elimination only - starts the bottom half of seeds in the losers bracket
	SplitParticipants *bool `json:"split_participants"`
}

type FreeForAllOptions struct {
	MaxParticipants *int64 `json:"max_participants"`
}

type GroupStageOptions struct {
	GroupSize                         *int64 `json:"group_size"`
	ParticipantCountToAdvancePerGroup *int64 `json:"participant_count_to_advance_per_group"`
	// Round robin format only, determines the primary ranking stat.
	RankedBy *RankedByEnum `json:"ranked_by"`
	// Round robin format only, where 1 = single round robin, 2 = double round robin, 3 = triple
	// round robin
	RrIterations *int64 `json:"rr_iterations"`
	// For 'custom' ranked by round robin group stage
	RrPtsForGameTie *float64 `json:"rr_pts_for_game_tie"`
	// For 'custom' ranked by round robin group stage
	RrPtsForGameWin *float64 `json:"rr_pts_for_game_win"`
	// For 'custom' ranked by round robin group stage
	RrPtsForMatchTie *float64 `json:"rr_pts_for_match_tie"`
	// For 'custom' ranked by round robin group stage
	RrPtsForMatchWin *float64 `json:"rr_pts_for_match_win"`
	// Double elimination only - starts the bottom half of seeds in the losers bracket
	SplitParticipants *bool      `json:"split_participants"`
	StageType         *StageType `json:"stage_type"`
}

type MatchOptions struct {
	// Whether or not to allow match attachment uploads
	AcceptAttachments *bool `json:"accept_attachments"`
	// For single or double elimination only, consolation matches will be run to break ties up
	// to this final placement.
	ConsolationMatchesTargetRank *int64 `json:"consolation_matches_target_rank"`
}

type Notifications struct {
	UponMatchesOpen    *bool `json:"upon_matches_open"`
	UponTournamentEnds *bool `json:"upon_tournament_ends"`
}

type RegistrationOptions struct {
	// Number of minutes for check-in prior to tournament start time. Must be a multiple of 5.
	CheckInDuration *int64 `json:"check_in_duration"`
	// Allow registered Challonge users to self-register for this tournament
	OpenSignup *bool `json:"open_signup"`
	// Maximum number of participants allowed in the tournament before the waitlist kicks in
	SignupCap *int64 `json:"signup_cap"`
}

type RoundRobinOptions struct {
	Iterations     *int64        `json:"iterations"`
	PtsForGameTie  *float64      `json:"pts_for_game_tie"`
	PtsForGameWin  *float64      `json:"pts_for_game_win"`
	PtsForMatchTie *float64      `json:"pts_for_match_tie"`
	PtsForMatchWin *float64      `json:"pts_for_match_win"`
	Ranking        *RankedByEnum `json:"ranking"`
}

type SeedingOptions struct {
	HideSeeds *bool `json:"hide_seeds"`
	// When true, seeding rules are ignored and participants are placed in the bracket from top
	// to bottom.
	SequentialPairings *bool `json:"sequential_pairings"`
}

type StationOptions struct {
	// Automatically assign stations to playable matches (requires one or more stations)
	AutoAssign *bool `json:"auto_assign"`
	// When true, playable matches won't start until they have a station assigned to them
	OnlyStartMatchesWithAssignedStations *bool `json:"only_start_matches_with_assigned_stations"`
}

type SwissOptions struct {
	PtsForGameTie  *float64 `json:"pts_for_game_tie"`
	PtsForGameWin  *float64 `json:"pts_for_game_win"`
	PtsForMatchTie *float64 `json:"pts_for_match_tie"`
	PtsForMatchWin *float64 `json:"pts_for_match_win"`
	Rounds         *int64   `json:"rounds"`
}

// When blank, the losers bracket winner has to beat the winners bracket winner twice.
type GrandFinalsModifierEnum string

const (
	GrandFinalsModifier GrandFinalsModifierEnum = ""
	SingleMatch         GrandFinalsModifierEnum = "single match"
	Skip                GrandFinalsModifierEnum = "skip"
)

// Round robin format only, determines the primary ranking stat.
type RankedByEnum string

const (
	Custom            RankedByEnum = "custom"
	GameWINS          RankedByEnum = "game wins"
	GameWinPercentage RankedByEnum = "game win percentage"
	MatchWINS         RankedByEnum = "match wins"
	PointsDifference  RankedByEnum = "points difference"
	PointsScored      RankedByEnum = "points scored"
	Rank              RankedByEnum = ""
)

type StageType string

const (
	StageTypeDoubleElimination StageType = "double elimination"
	StageTypeRoundRobin        StageType = "round robin"
	StageTypeSingleElimination StageType = "single elimination"
)

type TournamentType string

const (
	FreeForAll                      TournamentType = "free for all"
	Swiss                           TournamentType = "swiss"
	TournamentTypeDoubleElimination TournamentType = "double elimination"
	TournamentTypeRoundRobin        TournamentType = "round robin"
	TournamentTypeSingleElimination TournamentType = "single elimination"
)

func (c *Client) GetTournament(ctx context.Context, tourneySlug string) (Tournament, error) {
	slug := ExtractSlug(tourneySlug)

	results, err := GetData[ApidogTournamentResponse](c, ctx, getTournament, slug)
	if err != nil {
		return Tournament{}, fmt.Errorf("getTournament | Error getting data: %w", err)
	}

	if results == nil {
		return Tournament{}, fmt.Errorf("getTournament | API returned nil result for slug: %s", slug)
	}

	if results.Data == nil {
		return Tournament{}, fmt.Errorf("getTournament | Response 'data' field is missing for slug: %s (check API version or slug)", slug)
	}
	if results.Data.Attributes == nil {
		return Tournament{}, fmt.Errorf("getTournament | Tournament 'attributes' are missing for slug: %s", slug)
	}

	return *results.Data.Attributes, nil
}
