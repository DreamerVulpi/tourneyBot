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

	getPagesCount = `
	query getPagesCount($phaseGroupId: ID!, $states: [Int]){
		phaseGroup(id:$phaseGroupId){
			id
			sets (
				filters: {state: $states}
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

	GetPhaseGroupSets = `
	query getSets($phaseGroupId: ID!, $page:Int!, $perPage:Int!, $states: [Int]){
		phaseGroup(id:$phaseGroupId){
			id
			sets(
				page: $page
				perPage: $perPage
				sortType: STANDARD
				filters: {state: $states}
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
