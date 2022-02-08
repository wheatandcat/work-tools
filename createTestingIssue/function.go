package createTestingIssue

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"gocloud.dev/pubsub"
)

type IssueConfig struct {
	Repository string
	StartText  string
	EndText    string
}

func CreateTestingIssuePub(ctx context.Context, m *pubsub.Message) error {
	InstallationID, err := strconv.Atoi(os.Getenv("INSTALLATION_ID"))
	if err != nil {
		return err
	}
	githubAppID := os.Getenv("GITHUB_APP_ID")

	appID, err := strconv.ParseInt(githubAppID, 10, 64)
	if err != nil {

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
		StartText:  "## テスト内容",
		EndText:    "## テスト内容（エンジニア向け）",
	}}

	iss := []Issue{}

	for _, item := range items {
		r := Request{
			Token:      token,
			Owner:      os.Getenv("GITHUB_OWNER"),
			Repository: item.Repository,
			StartText:  item.StartText,
			EndText:    item.EndText,
		}

		tiss, err := CreateIssues(r)
		if err != nil {
			return err
		}

		iss = append(iss, tiss...)
	}

	if len(iss) == 0 {
		// issueがないので終了
		return nil
	}

	err = PostSlack(iss)
	if err != nil {
		return err
	}

	return nil
}
