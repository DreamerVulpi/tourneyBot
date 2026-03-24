package startgg

import (
	"github.com/dreamervulpi/tourneyBot/internal/entity/startgg"
)

func (c *Client) GetPagesCount(phaseGroupID int64, states []int) (int, error) {
	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
		"states":       states,
	}

	results, err := GetData[startgg.RawPagesDataCount](c, startgg.GetPagesCount, variables)
	if err != nil {
		return 0, err
	}

	return results.Data.PhaseGroup.Sets.PageInfo.Total, nil
}
