package startgg

const (
	getEvent = `
	query getEvent($slug: String) {
		event(slug: $slug) {
			id
			name
		}
	},
	`
	getPhaseGroupState = `
	query PhaseGroupState($phaseGroupId: ID!){
		phaseGroup(id:$phaseGroupId){
			id
			state
		}
	}`

	getPhaseGroupSets = `
	query PhaseGroupSets($phaseGroupId: ID!, $page:Int!, $perPage:Int!){
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
