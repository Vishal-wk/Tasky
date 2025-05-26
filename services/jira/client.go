package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

)

const (
	startAt = 0
	maxResults = 50
)
type client struct {
	domain string
}

func (c *client) projects() Project {
	projects := []map[string]interface{}{}
	
	for {
		url := fmt.Sprintf("%s/rest/api/3/project/search?startAt=%d&maxResults=%d", c.domain, startAt, maxResults)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}
		auth := base64.StdEncoding.EncodeToString([]byte(email + ":" + apiToken))
		req.Header.Set("Authorization", "Basic "+auth)
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Fatalf("Jira API returned status %d", resp.StatusCode)
		}

		var result struct {
			IsLast     bool                     `json:"isLast"`
			StartAt    int                      `json:"startAt"`
			MaxResults int                      `json:"maxResults"`
			Total      int                      `json:"total"`
			Values     []map[string]interface{} `json:"values"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			log.Fatalf("Failed to parse project response: %v", err)
		}

		projects = append(projects, result.Values...)
		if result.IsLast || len(result.Values) == 0 {
			break
		}
		startAt += maxResults
	}
	return projects
}
