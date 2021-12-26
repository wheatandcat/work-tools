package postAssignableIssues

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/bradleyfalzon/ghinstallation"
)

type IssueConfig struct {
	Repository string
	Categories []string
}

func PostAssignableIssuesPubSub(ctx context.Context, m *pubsub.Message) error {
	InstallationID, err := strconv.Atoi(os.Getenv("INSTALLATION_ID"))
	if err != nil {
		return err
	}
	githubAppID := os.Getenv("GITHUB_APP_ID")

	appID, err := strconv.ParseInt(githubAppID, 10, 64)
	if err != nil {
		return err
	}
	key := os.Getenv("GITHUB_APP_PRIVATE_KEY")
	tr := http.DefaultTransport
	itr, err := ghinstallation.New(tr, appID, int64(InstallationID), []byte(key))
	if err != nil {
		return err
	}

	token, err := itr.Token(ctx)

	if err != nil {
		return err
	}

	items := []IssueConfig{{
		Repository: "tool-test",
		Categories: []string{"RFP", "ドキュメンテーション", "設計前", "運用改善"},
	}}
	text := ""

	for _, item := range items {
		r := Request{
			Token:      token,
			Owner:      os.Getenv("GITHUB_OWNER"),
			Repository: item.Repository,
			Categories: item.Categories,
		}

		t, err := GetIssueText(r)
		if err != nil {
			return err
		}
		text += t
	}

	rn := RequestNote{
		Token:   os.Getenv("NOTE_TOKEN"),
		ID:      os.Getenv("NOTE_ID"),
		Host:    os.Getenv("NOTE_HOST"),
		Content: text,
	}

	if err := PostNote(rn); err != nil {
		return err
	}

	return nil
}
