package ghclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) query(q string, variables map[string]any, result any) error {
	token, err := c.bearerToken()
	if err != nil {
		return err
	}

	body, _ := json.Marshal(graphqlRequest{Query: q, Variables: variables})
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("graphql request failed: %s", resp.Status)
	}

	var gqlResp graphqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return err
	}
	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}

	return json.Unmarshal(gqlResp.Data, result)
}
