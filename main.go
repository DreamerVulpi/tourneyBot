package main

import (
	"fmt"

	"github.com/dreamervulpi/tourneybot/internal/config"
	"github.com/dreamervulpi/tourneybot/internal/startgg"
)

func main() {
	cfg, err := config.LoadConfig("internal/config/config.toml")
	if err != nil {
		fmt.Println("Not loaded configation")
	}

	startgg.AuthToken = cfg.Startgg.Token
	rawPhaseGroupData, err := startgg.GetPhaseGroupSets(2553660, 1, 60)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("PhaseGroupID: ", rawPhaseGroupData.Data.PhaseGroup.Id)
	fmt.Println("Sets in group: ", rawPhaseGroupData.Data.PhaseGroup.Sets.PageInfo)
	for _, set := range rawPhaseGroupData.Data.PhaseGroup.Sets.Nodes {
		fmt.Println("set ID: ", set.Id)
		fmt.Println("State set: ", set.State)
		fmt.Println("player 1 | ID: ", set.Slots[0].Entrant.Id)
		fmt.Println("player 1 | Nickname: ", set.Slots[0].Entrant.Participants[0].GamerTag)
		fmt.Println("player 1 | TEKKEN ID ", set.Slots[0].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID)
		fmt.Println("player 1 | Discord: ", set.Slots[0].Entrant.Participants[0].User.Authorizations[0].Discord)
		fmt.Println("player 2 | ID: ", set.Slots[1].Entrant.Id)
		fmt.Println("player 2 | Nickname: ", set.Slots[1].Entrant.Participants[0].GamerTag)
		fmt.Println("player 2 | TEKKEN ID ", set.Slots[1].Entrant.Participants[0].ConnectedAccounts.Tekken.TekkenID)
		fmt.Println("player 2 | Discord: ", set.Slots[1].Entrant.Participants[0].User.Authorizations[0].Discord)
	}
}
