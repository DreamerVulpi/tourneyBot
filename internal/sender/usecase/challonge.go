package usecase

import (
	"context"
	"log"
	"strconv"

	"fmt"

	"github.com/dreamervulpi/tourneyBot/challonge"
	"github.com/dreamervulpi/tourneyBot/internal/sender"
)

type ChallongeMatchAdapter struct {
	Client         *challonge.Client
	TournamentSlug string
	DebugMode      bool
	TestUser       sender.Participant
}

func (c ChallongeMatchAdapter) GetPlatformTournamentName() string {
	return "challonge"
}

func (c ChallongeMatchAdapter) GetTournamentSlug() string {
	return c.TournamentSlug
}

func (c ChallongeMatchAdapter) GetSetsData(ctx context.Context) ([]sender.SetData, error) {
	// TODO: Complete method
	// https://challonge.com/ru/tournamentdciii

	tournament, err := c.Client.GetTournament(ctx, c.TournamentSlug)
	if err != nil {
		return nil, fmt.Errorf("GetSetsData | Challonge | get tournament error: %w", err)
	}

	states := []challonge.State{challonge.Open}
	if c.DebugMode {
		states = []challonge.State{challonge.Open, challonge.Pending, challonge.Complete}
	}

	matches, err := c.Client.GetMatches(ctx, c.TournamentSlug, states)
	if err != nil {
		return nil, fmt.Errorf("GetSetsData | Challonge | Can't get data of sets: %v", err)
	}

	var matchesData []sender.SetData

	for _, match := range matches {
		if ctx.Err() != nil {
			break
		}

		// TODO: ConvertData for ChallongeContacts
		rawP1, err := c.Client.GetParticipant(ctx, c.TournamentSlug, match.PointsByParticipant[0].ParticipantID)
		if err != nil {
			log.Printf("GetSetsData | Challonge | Can't get data of player 1 (%v) from match (%v): %v\n", match.PointsByParticipant[0].ParticipantID, match.Identifier, err)
		}
		rawP2, err := c.Client.GetParticipant(ctx, c.TournamentSlug, match.PointsByParticipant[1].ParticipantID)
		if err != nil {
			log.Printf("GetSetsData | Challonge | Can't get data of player 1 (%v) from match (%v): %v\n", match.PointsByParticipant[0].ParticipantID, match.Identifier, err)
		}

		p1 := c.ConvertContacts(rawP1)
		p2 := c.ConvertContacts(rawP2)

		// TODO: isFinals | how get data about rounds? all
		isFinals := false

		var matchID int64
		convertedMatchID, err := strconv.ParseInt(match.ID, 10, 64)
		if err != nil || convertedMatchID <= 0 {
			log.Printf("GetSetsData | Challonge | Can't convert match id type string (%v) to int64: %v\n", match.ID, err)
			matchID = 0
		} else {
			matchID = convertedMatchID
		}
		set := sender.SetData{
			TournamentName: tournament.Name,
			SetID:          matchID,
			// TODO: StreamName, StreamSourse
			RoundNum:       match.Round,
			PhaseGroupId:   0,
			IsFinals:       isFinals,
			ContactPlayer1: p1,
			ContactPlayer2: p2,
			FullInviteLink: fmt.Sprintf("https://challonge.com/en/matches/%v/chat.html", matchID),
		}
		matchesData = append(matchesData, set)
	}

	return matchesData, nil
}

func (c *ChallongeMatchAdapter) ConvertContacts(data challonge.ParticipantOutput) sender.Participant {
	p := sender.Participant{
		MessenagerLogin: "N/D",
		GameID:          "N/D",
		GameNickname:    "N/D",
	}

	if len(data.Name) != 0 {
		p.GameNickname = data.Name
	} else if len(data.Username) != 0 {
		p.GameNickname = data.Username
	}

	return p
}
