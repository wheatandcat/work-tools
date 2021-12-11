package createissue

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
)

const RepoOwner = "wheatandcat"
const Repo = "gas-tools"
const IssueNumber = 17

func CreateIssue(w http.ResponseWriter, r *http.Request) {
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

	client := github.NewClient(&http.Client{
		Transport: itr,
		Timeout:   5 * time.Second,
	})

	ctx := context.Background()

	body := "hello"
	comment := &github.IssueComment{
		Body: &body,
	}
	if _, _, err := client.Issues.CreateComment(ctx, RepoOwner, Repo, IssueNumber, comment); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
