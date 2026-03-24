package startgg

import (
	"github.com/dreamervulpi/tourneyBot/internal/entity/startgg"
)

func (c *Client) GetListGroups(slug string) ([]startgg.PhaseGroupId, error) {
	var variables = map[string]any{
		"slug": slug,
	}

	results, err := GetData[startgg.RawPhaseGroupsData](c, startgg.GetListPhaseGroups, variables)
	if err != nil {
		return nil, err
	}

	return results.Data.Event.PhaseGroups, nil
}
