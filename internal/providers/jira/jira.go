package jira

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// JiraClient holds authentication and base URL
type JiraClient struct {
    BaseURL  string
    Email    string // required for API token auth
    APIToken string
}

// JiraIssue is a minimal struct for tasks/stories
type JiraIssue struct {
    Key    string `json:"key,omitempty"`
    Fields struct {
        Summary     string `json:"summary"`
        Description string `json:"description,omitempty"`
    } `json:"fields"`
}

// NewJiraClient creates a new client instance
func NewJiraClient(baseURL, email, apiToken string) *JiraClient {
    return &JiraClient{
        BaseURL:  baseURL,
        Email:    email,
        APIToken: apiToken,
    }
}

// GetIssue fetches an issue by key
func (c *JiraClient) GetIssue(issueKey string) (*JiraIssue, error) {
    url := fmt.Sprintf("%s/rest/api/3/issue/%s", c.BaseURL, issueKey)

    req, _ := http.NewRequest("GET", url, nil)
    req.SetBasicAuth(c.Email, c.APIToken)
    req.Header.Set("Accept", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("jira error: %s", string(body))
    }

    var issue JiraIssue
    if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
        return nil, err
    }
    return &issue, nil
}

// CreateIssue creates a new issue
func (c *JiraClient) CreateIssue(projectKey, summary, description string) (*JiraIssue, error) {
    url := fmt.Sprintf("%s/rest/api/3/issue", c.BaseURL)

    payload := map[string]interface{}{
        "fields": map[string]interface{}{
            "project": map[string]string{"key": projectKey},
            "summary": summary,
            "description": map[string]string{
                "type": "doc",
                "version": "1",
                "content": description,
            },
            "issuetype": map[string]string{"name": "Task"},
        },
    }
    data, _ := json.Marshal(payload)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
    req.SetBasicAuth(c.Email, c.APIToken)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 201 {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("jira error: %s", string(body))
    }

    var issue JiraIssue
    if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
        return nil, err
    }
    return &issue, nil
}
