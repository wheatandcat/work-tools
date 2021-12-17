package createissue

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
)

type Response struct {
	Issues []Issue `json:"issues"`
}

func CreateIssue(w http.ResponseWriter, r *http.Request) {
	at := r.Header.Get("Authorization")
	if at != os.Getenv("VERIFY_ID_TOKEN") {
		http.Error(w, fmt.Errorf("error: verifying ID token, %s", at).Error(), http.StatusInternalServerError)
		return
	}

	var d struct {
		ID           int      `json:"id"`
		Priority     string   `json:"priority"`
		Title        string   `json:"title"`
		Body         string   `json:"body"`
		Env          string   `json:"env"`
		Image        string   `json:"image"`
		Version      string   `json:"version"`
		Repositories []string `json:"repositories"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	InstallationID, err := strconv.Atoi(os.Getenv("INSTALLATION_ID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	githubAppID := os.Getenv("GITHUB_APP_ID")

	appID, err := strconv.ParseInt(githubAppID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := os.Getenv("GITHUB_APP_PRIVATE_KEY")
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, appID, int64(InstallationID), []byte(key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	token, err := itr.Token(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{}

	for _, repository := range d.Repositories {
		req := Request{
			Token:      token,
			Owner:      os.Getenv("GITHUB_OWNER"),
			Version:    d.Version,
			Repository: repository,
			ID:         d.ID,
			Priority:   d.Priority,
			Title:      d.Title,
			Body:       d.Body,
			Env:        d.Env,
			Image:      d.Image,
		}
		r, err := Create(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		is := Issue{
			ID:    strconv.Itoa(r.ID),
			Title: r.CreateIssue.Issue.Title,
			URL:   r.CreateIssue.Issue.URL,
		}
		res.Issues = append(res.Issues, is)
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
