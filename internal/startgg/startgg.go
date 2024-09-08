package startgg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type State int

const (
	IsNotStarted State = 1
	InProcess    State = 2
	IsDone       State = 3
)

type StateEvent string

const (
	Created   StateEvent = "CREATED"
	Active    StateEvent = "ACTIVE"
	Completed StateEvent = "COMPLETED"
)

type Client struct {
	AuthToken string
	Client    *http.Client
}

func NewClient(token string, clt *http.Client) *Client {
	return &Client{
		AuthToken: token,
		Client:    clt,
	}
}

func PrepareQuery(query string, variables map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
}

func (c *Client) RunQuery(query []byte) ([]byte, error) {
	// Creates the POST request and checks for errors.
	req, err := http.NewRequest("POST", "https://api.start.gg/gql/alpha", bytes.NewBuffer(query))
	if err != nil {
		return nil, errors.Join(errors.New("HTTP Request - "), err)
	}

	// Sets the headers within the request.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))

	// Sends the request and receives the response of it.
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("HTTP Response - "), err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("read Data - "), err)
	}

	validation, err := validateData(data)
	if err != nil {
		return nil, err
	}
	if validation != "" {
		return nil, errors.Join(errors.New("data Validation - "), err)
	}

	return data, nil
}

func validateData(data []byte) (string, error) {
	results := &FailedCall{}

	err := json.Unmarshal(data, results)
	if err != nil {
		return "", errors.New("unable To Validate Data")
	}

	return results.Message, nil
}
