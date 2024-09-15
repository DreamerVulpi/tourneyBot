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
	getTournament = `
	query getTournament($tourneySlug:String!) {
		tournament(slug: $tourneySlug) {
			id
			name
			state
		}
	}
	`
	// TODO: Set filter to 1
	getPagesCount = `
	query getPagesCount($phaseGroupId: ID!){
		phaseGroup(id:$phaseGroupId){
			id
			sets (
				filters: {state: 3}
			){
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
	// TODO: Set filter to 1
	getPhaseGroupSets = `
	query getSets($phaseGroupId: ID!, $page:Int!, $perPage:Int!){
		phaseGroup(id:$phaseGroupId){
			id
			sets(
				page: $page
				perPage: $perPage
				sortType: STANDARD
				filters: {state: 3}
			){
			pageInfo{
				total
			}
			nodes{
					id
					state
					stream {
						streamName
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
