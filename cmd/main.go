package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

var oauthConfig = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  redirectURL,
	Scopes:       []string{"read:jira-user", "read:jira-work"}, // adjust scopes as needed
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://auth.atlassian.com/authorize",
		TokenURL: "https://auth.atlassian.com/oauth/token",
	},
}

func main() {
	// Start HTTP server to handle OAuth callback
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/callback", handleCallback)

	fmt.Println("Starting server at http://localhost:8080 ...")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Open the user's browser for authorization
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Please open the following URL in your browser to authorize:\n%v\n", authURL)

	// Wait forever - the callback handler will exit after token received
	select {}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Jira OAuth 2.0 3LO example. Navigate to /callback after authorization.")
}
func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Authorization successful! You can close this window.\n")

	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Refresh Token:", token.RefreshToken)
	fmt.Println("Token Expiry:", token.Expiry.Format(time.RFC1123))

	client := oauthConfig.Client(ctx, token)

	// Step 1: Get account info (hosted API)
	resp, err := client.Get("https://api.atlassian.com/me")
	if err != nil {
		log.Fatalf("Failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Fatalf("Failed to decode user info: %v", err)
	}
	fmt.Println("\n🔐 User Info from Atlassian API:")
	printJSON(userInfo)

	// Step 2: Get all Jira projects from your Jira domain
	jiraAPI := fmt.Sprintf("%s/rest/api/3/project", jiraDomain)
	resp, err = client.Get(jiraAPI)
	if err != nil {
		log.Fatalf("Failed to get Jira projects: %v", err)
	}
	defer resp.Body.Close()

	var projects interface{}
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		log.Fatalf("Failed to decode Jira projects: %v", err)
	}
	fmt.Println("\n📁 Jira Projects from", jiraDomain+":")
	printJSON(projects)

	os.Exit(0)
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
