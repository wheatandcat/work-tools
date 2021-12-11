package createissue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
)

const RepoOwner = "wheatandcat"
const Repo = "gas-tools"
const IssueNumber = 17

func CreateIssue(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprint(w, "Hello, World!")
		return
	}

	InstallationID, err := strconv.Atoi(os.Getenv("INSTALLATION_ID"))
	if err != nil {
		log.Fatal(err)
	}
	githubAppID := os.Getenv("GITHUB_APP_ID")

	appID, err := strconv.ParseInt(githubAppID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	key := os.Getenv("GITHUB_APP_PRIVATE_KEY")
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, appID, int64(InstallationID), []byte(key))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	token, err := itr.Token(ctx)

	if err != nil {
		log.Fatal(err)
	}

	for _, repository := range d.Repositories {
		req := Request{
			Token:      token,
			Owner:      RepoOwner,
			Version:    d.Version,
			Repository: repository,
			ID:         d.ID,
			Priority:   d.Priority,
			Title:      d.Title,
			Body:       d.Body,
			Env:        d.Env,
			Image:      d.Image,
		}
		if err := Create(req); err != nil {
			log.Fatal(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
