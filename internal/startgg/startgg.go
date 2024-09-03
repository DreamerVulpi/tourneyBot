package startgg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	AuthToken string
	client    = &http.Client{
		Timeout: time.Second * 10,
	}
)

func SetAuthToken(token string) {
	AuthToken = token
}

func Token() bool {
	return len(AuthToken) > 0
}

func PrepareQuery(query string, variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
}

func RunQuery(query []byte) ([]byte, error) {
	// Creates the POST request and checks for errors.
	req, err := http.NewRequest("POST", "https://api.start.gg/gql/alpha", bytes.NewBuffer(query))
	if err != nil {
		return nil, fmt.Errorf("HTTP Request - %w", err)
	}

	// Sets the headers within the request.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AuthToken))

	// Sends the request and receives the response of it.
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP Response - %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read data - %w", err)
	}

	validation := validateData(data)
	if validation != "" {
		return nil, fmt.Errorf("data validation - %s", validation)
	}

	return data, nil
}

func validateData(data []byte) string {
	results := &FailedCall{}

	err := json.Unmarshal(data, results)
	if err != nil {
		return "Unable To Validate Data"
	}

	return results.Message
}
