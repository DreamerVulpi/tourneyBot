package startgg

const (
	getListPhaseGroups = `
	query getListPhaseGroups($slug: String) {
		event(slug: $slug) {
			id
			name
			phaseGroups {
				id
			}
		}
	},
	`
	getPagesCount = `
	query getPagesCount($phaseGroupId: ID!){
		phaseGroup(id:$phaseGroupId){
			id
			sets {
				pageInfo{
					total
				}
			}
		}
	}
		`

	// TODO: WINNER FUNCTION
	getWinner = `
	query getWinner($setId: ID!){
		set(id:$setId) {
			winnerId
		}
	}
	`

	getPhaseGroupState = `
	query getPhaseGroupState($phaseGroupId: ID!){
		phaseGroup(id:$phaseGroupId){
			id
			state
		}
	}`

	getPhaseGroupSets = `
	query getSets($phaseGroupId: ID!, $page:Int!, $perPage:Int!){
		phaseGroup(id:$phaseGroupId){
			id
			sets(
				page: $page
				perPage: $perPage
				sortType: STANDARD
			){
			pageInfo{
				total
			}
			nodes{
					id
					state
					stream {
						streamSource
					}
					slots{
						entrant{
							id
							participants {
								gamerTag
								connectedAccounts
								user {
									authorizations(types: DISCORD) {
										externalUsername
									}
								}
							}
						}
					}
				}
			}
		}
	}`
)
