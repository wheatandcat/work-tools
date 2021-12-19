package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
)

func main() {
	InstallationID, err := strconv.Atoi(os.Getenv("INSTALLATION_ID"))
	if err != nil {
		panic(err)
	}
	githubAppID := os.Getenv("GITHUB_APP_ID")

	appID, err := strconv.ParseInt(githubAppID, 10, 64)
	if err != nil {
		panic(err)
	}
	key := os.Getenv("GITHUB_APP_PRIVATE_KEY")
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, appID, int64(InstallationID), []byte(key))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	token, err := itr.Token(ctx)

	if err != nil {
		panic(err)
	}

	repositories := []string{"gas-tools"}
	text := ""

	for _, rep := range repositories {
		r := Request{
			Token:      token,
			Owner:      os.Getenv("GITHUB_OWNER"),
			Repository: rep,
		}

		t, err := GetIssueText(r)
		if err != nil {
			panic(err)
		}
		text += t
	}

	log.Println(text)

	rn := RequestNote{
		Token:   os.Getenv("NOTE_TOKEN"),
		ID:      os.Getenv("NOTE_ID"),
		Host:    os.Getenv("NOTE_HOST"),
		Content: text,
	}

	if err := PostNote(rn); err != nil {
		panic(err)
	}
}
