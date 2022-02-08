package createTestingIssue

import (
	"os"

	"github.com/slack-go/slack"
)

func PostSlack(iss []Issue) error {
	text := "<!here> \n"
	text += "以下のテストが作成されました。ご確認お願いします 🙏 \n"

	for _, is := range iss {
		text += " • <" + is.URL + "|" + is.Title + ">\n"
	}

	tkn := os.Getenv("SLACK_TOKEN")
	c := slack.New(tkn)
	_, _, err := c.PostMessage(os.Getenv("SLACK_CHANNEL"), slack.MsgOptionText(text, false), slack.MsgOptionDisableMediaUnfurl(), slack.MsgOptionDisableLinkUnfurl())
	if err != nil {
		return err
	}

	return nil
}
