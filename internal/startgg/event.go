package startgg

import (
	"encoding/json"
	"errors"
	"fmt"
)

type RawEventData struct {
	Data   DataEvent `json:"data"`
	Errors []Errors  `json:"errors"`
}

type DataEvent struct {
	Event Event `json:"event"`
}

type Event struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func GetEvent(slug string) (*RawEventData, error) {
	if !token() {
		return &RawEventData{}, errors.New("token verification - authentication token not set")
	}

	var variables = map[string]any{
		"slug": slug,
	}

	query, err := json.Marshal(prepareQuery(getEvent, variables))
	if err != nil {
		return &RawEventData{}, fmt.Errorf("JSON Marshal - %w", err)
	}

	data, err := runQuery(query)
	if err != nil {
		return &RawEventData{}, err
	}

	results := &RawEventData{}
	err = json.Unmarshal(data, results)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return results, nil
}
