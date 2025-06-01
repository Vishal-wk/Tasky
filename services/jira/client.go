package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Vishal/Tasky/config"

)
type Client struct{
	Email string
	APIToken string
	Domain string
}

func New(email, APIToken, domain string)Client{
 return Client{
	Email: email,
	APIToken: APIToken,
	Domain: domain,
 }
}
func GetAllProjects(cfg config.Config) []Value {
	var allProjects []Value
	url := fmt.Sprintf("%s/rest/api/3/project/search?maxResults=50", cfg.Domain)

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		req.SetBasicAuth(cfg.Email, cfg.APIToken)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			log.Fatalf("Jira API returned status %d: %s", resp.StatusCode, string(body))
		}

		var projectResp Project
		if err := json.NewDecoder(resp.Body).Decode(&projectResp); err != nil {
			log.Fatalf("Failed to decode response: %v", err)
		}

		allProjects = append(allProjects, projectResp.Values...)

		if projectResp.IsLast || projectResp.NextPage == "" {
			break
		}

		url = projectResp.NextPage
	}

	return allProjects
}

func GetIssues(cfg config.Config, projectKey string) []map[string]interface{} {
	var allIssues []map[string]interface{}
	startAt := 0
	maxResults := 50

	for {
		jql := fmt.Sprintf("project=%s", projectKey)
		url := fmt.Sprintf("%s/rest/api/3/search?jql=%s&startAt=%d&maxResults=%d", cfg.Domain, jql, startAt, maxResults)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		auth := basicAuth(cfg.Email, cfg.APIToken)
		req.Header.Set("Authorization", auth)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			log.Fatalf("Jira API returned status %d: %s", resp.StatusCode, string(body))
		}

		var data struct {
			Issues     []map[string]interface{} `json:"issues"`
			Total      int                      `json:"total"`
			StartAt    int                      `json:"startAt"`
			MaxResults int                      `json:"maxResults"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatalf("Failed to decode issues: %v", err)
		}

		allIssues = append(allIssues, data.Issues...)

		if len(data.Issues) < maxResults {
			break
		}
		startAt += maxResults
	}

	return allIssues
}

func basicAuth(email, token string) string {
	creds := fmt.Sprintf("%s:%s", email, token)
	return "Basic " + base64Encode(creds)
}

func base64Encode(s string) string {
	return strings.TrimRight(strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(string([]byte(s))), "=")
}

func ValidateToken(token string) (string, error) {
    req, err := http.NewRequest("GET", "https://your-domain.atlassian.net/rest/api/3/myself", nil)
    if err != nil {
        return "", err
    }

    req.Header.Add("Authorization", "Bearer "+token)
    req.Header.Add("Accept", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 200 {
        return "Token is valid", nil
    }

    return "", fmt.Errorf("Invalid token: status %d", resp.StatusCode)
}