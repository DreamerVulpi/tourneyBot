package startgg

import (
	"log"
)

type RawPagesDataCount struct {
	Data   DataPhaseGroupSets `json:"data"`
	Errors []Errors           `json:"errors"`
}

func (c *Client) GetPagesCount(phaseGroupID int64, states []int) (int, error) {
	var variables = map[string]any{
		"phaseGroupId": phaseGroupID,
		"states":       states,
	}
	log.Printf("SENDING VARS: %+v", variables)
	results, err := GetData[RawPagesDataCount](c, getPagesCount, variables)
	if err != nil {
		return 0, err
	}

	return results.Data.PhaseGroup.Sets.PageInfo.Total, nil
}
