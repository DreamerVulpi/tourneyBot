package startgg

import (
	"encoding/json"
	"fmt"
	"testing"

	"os"

	"github.com/dreamervulpi/tourneyBot/startgg"
	"github.com/stretchr/testify/assert"
)

// Test function GetData using request TestGetPhaseGroupSets
func TestGetData(t *testing.T) {
	variables := map[string]any{
		"phaseGroupId": 2562129,
		"page":         1,
		"perPage":      4,
	}

	_, err := json.Marshal(startgg.PrepareQuery(startgg.TestGetPhaseGroupSets, variables))
	if err != nil {
		assert.Error(t, fmt.Errorf("JSON Marshal - %w", err))
	}

	rawData, err := os.ReadFile("rawDataSets.json")
	if err != nil {
		assert.Error(t, fmt.Errorf("RunQuery - %w", err))
	}

	expected := &startgg.RawPhaseGroupSetsData{
		Data: startgg.DataPhaseGroupSets{
			PhaseGroup: startgg.PhaseGroupSets{
				Id: 2562129,
				Sets: startgg.Sets{
					PageInfo: startgg.PageInfo{Total: 9},
					Nodes: []startgg.Nodes{
						{
							Id:            77972366,
							State:         3,
							FullRoundText: "Winners Final",
							Round:         3,
							Slots: []startgg.Slots{
								{
									Entrant: startgg.Entrant{
										Id: 17340928,
										Participants: []startgg.Participants{
											{
												GamerTag: "Pizzduk",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{
														TekkenID: "636h-F3rQ-Hqd2",
													},
												},
												User: startgg.User{
													Authorizations: []startgg.Authorizations{
														{
															Discord: "pizzduk",
														},
													},
												},
											},
										},
									},
								},
								{
									Entrant: startgg.Entrant{
										Id: 17299461,
										Participants: []startgg.Participants{
											{
												GamerTag: "Mr_Shadow_",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{},
												},
												User: startgg.User{
													Authorizations: []startgg.Authorizations{
														{
															Discord: "mr_shadow_",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Id:            77972382,
							State:         3,
							FullRoundText: "Losers Final",
							Round:         -6,
							Slots: []startgg.Slots{
								{
									Entrant: startgg.Entrant{
										Id: 17340928,
										Participants: []startgg.Participants{
											{
												GamerTag: "Pizzduk",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{
														TekkenID: "636h-F3rQ-Hqd2",
													},
												},
												User: startgg.User{
													Authorizations: []startgg.Authorizations{
														{
															Discord: "pizzduk",
														},
													},
												},
											},
										},
									},
								},
								{
									Entrant: startgg.Entrant{
										Id: 17373189,
										Participants: []startgg.Participants{
											{
												GamerTag: "cleverdemon",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{},
												},
												User: startgg.User{
													Authorizations: nil,
												},
											},
										},
									},
								},
							},
						},
						{
							Id:            77972381,
							State:         3,
							FullRoundText: "Losers Semi-Final",
							Round:         -5,
							Slots: []startgg.Slots{
								{
									Entrant: startgg.Entrant{
										Id: 17327065,
										Participants: []startgg.Participants{
											{
												GamerTag: "AlexSouls",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{
														TekkenID: "6624-jt93-MtEE",
													},
												},
												User: startgg.User{
													Authorizations: []startgg.Authorizations{
														{
															Discord: "alexsouls",
														},
													},
												},
											},
										},
									},
								},
								{
									Entrant: startgg.Entrant{
										Id: 17373189,
										Participants: []startgg.Participants{
											{
												GamerTag: "cleverdemon",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{},
												},
												User: startgg.User{
													Authorizations: nil,
												},
											},
										},
									},
								},
							},
						},
						{
							Id:            77972364,
							State:         3,
							FullRoundText: "Winners Semi-Final",
							Round:         2,
							Slots: []startgg.Slots{
								{
									Entrant: startgg.Entrant{
										Id: 17373189,
										Participants: []startgg.Participants{
											{
												GamerTag: "cleverdemon",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{},
												},
												User: startgg.User{
													Authorizations: nil,
												},
											},
										},
									},
								},
								{
									Entrant: startgg.Entrant{
										Id: 17340928,
										Participants: []startgg.Participants{
											{
												GamerTag: "Pizzduk",
												ConnectedAccounts: startgg.ConnectedAccounts{
													Tekken: startgg.Tekken8{
														TekkenID: "636h-F3rQ-Hqd2",
													},
												},
												User: startgg.User{
													Authorizations: []startgg.Authorizations{
														{
															Discord: "pizzduk",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	actual := &startgg.RawPhaseGroupSetsData{}
	err = json.Unmarshal(rawData, actual)
	if err != nil {
		assert.Error(t, fmt.Errorf("JSON Unmarshal - %w", err))
	}

	assert.Equal(t, expected, actual)
}
