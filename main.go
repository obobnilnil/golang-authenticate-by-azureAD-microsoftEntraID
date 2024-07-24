package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var oauth2Config = &oauth2.Config{
	ClientID:     "your_client_id",
	ClientSecret: "your_client_secret",
	RedirectURL:  "your_callback_url", //  in my case RedirectURL:  "http://localhost:8080/callback"
	Scopes:       []string{"User.Read"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	},
}

func main() {
	http.HandleFunc("/", startLogin)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/logout", logout) // Add logout handler
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User Info: %s", string(data))
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Optional: Clear any local session or token storage
	// Redirect to Microsoft logout URL
	logoutURL := "https://login.microsoftonline.com/common/oauth2/v2.0/logout?post_logout_redirect_uri=http://localhost:8080"
	http.Redirect(w, r, logoutURL, http.StatusFound)
}
