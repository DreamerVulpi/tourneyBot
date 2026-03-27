package startgg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dreamervulpi/tourneyBot/internal/entity/startgg"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(clt *http.Client) *Client {
	return &Client{
		httpClient: clt,
	}
}

func PrepareQuery(query string, variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
}

func validateData(data []byte) (string, error) {
	results := &startgg.FailedCall{}

	err := json.Unmarshal(data, results)
	if err != nil {
		return "", errors.New("unable To Validate Data")
	}

	return results.Message, nil
}

// Execute query for get raw data
func (c *Client) RunQuery(query []byte) ([]byte, error) {
	// Creates the POST request and checks for errors.
	req, err := http.NewRequest("POST", "https://api.start.gg/gql/alpha", bytes.NewBuffer(query))
	if err != nil {
		return nil, errors.Join(errors.New("HTTP Request - "), err)
	}

	// Sets the headers within the request.
	req.Header.Set("Content-Type", "application/json")

	// Sends the request and receives the response of it.
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("HTTP Response - "), err)
	}
	defer res.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("read Data - "), err)
	}

	validation, err := validateData(data)
	if err != nil {
		return nil, err
	}
	if validation != "" {
		return nil, fmt.Errorf("data Validation: %s", validation)
	}

	return data, nil
}

// Execute query for get data from startgg according T type
func GetData[T any](c *Client, rawQuery string, variables map[string]any) (*T, error) {
	preparedQuery, err := json.Marshal(PrepareQuery(rawQuery, variables))
	if err != nil {
		return nil, fmt.Errorf("JSON Marshal - %w", err)
	}

	rawData, err := c.RunQuery(preparedQuery)
	if err != nil {
		return nil, fmt.Errorf("RunQuery - %w", err)
	}

	var results T
	err = json.Unmarshal(rawData, &results)
	if err != nil {
		return nil, fmt.Errorf("JSON Unmarshal - %w", err)
	}

	return &results, nil
}
