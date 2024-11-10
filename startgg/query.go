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
	// Test: Set filter to 3
	getPagesCount = `
	query getPagesCount($phaseGroupId: ID!){
		phaseGroup(id:$phaseGroupId){
			id
			sets (
				filters: {state: 1}
			){
				pageInfo{
					total
				}
			}
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
	// Test: Set filter to 3
	getPhaseGroupSets = `
	query getSets($phaseGroupId: ID!, $page:Int!, $perPage:Int!){
		phaseGroup(id:$phaseGroupId){
			id
			sets(
				page: $page
				perPage: $perPage
				sortType: STANDARD
				filters: {state: 1}
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
					fullRoundText
        			round
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
