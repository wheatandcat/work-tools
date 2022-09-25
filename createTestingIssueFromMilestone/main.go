package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/machinebox/graphql"
)

type C struct {
	GitHub GitHub
}

type GitHub struct {
	Token string `toml:"token"`
}

type IssueConfig struct {
	Repository string
	Owner      string
	Label      string
	StartText  string
	EndText    string
}

func main() {
	var config C
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	items := []IssueConfig{{
		Owner:      "wheatandcat",
		Repository: "work-tools",
		Label:      "テストを実施",
		StartText:  "## テスト内容",
		EndText:    "## テスト内容（エンジニア向け）",
	}}

	iss := []Issue{}

	trid, testMilestone, err := config.GitHub.getMilestoneNumber(os.Getenv("MILESTONE"), "wheatandcat", "tool-test")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		_, m, err := config.GitHub.getMilestoneNumber(os.Getenv("MILESTONE"), item.Owner, item.Repository)
		if err != nil {
			log.Fatal(err)
		}

		r := Request{
			Token:            config.GitHub.Token,
			Owner:            item.Owner,
			Milestone:        m,
			Label:            item.Label,
			TestMilestone:    testMilestone,
			TestRepositoryID: trid,
			Repository:       item.Repository,
			StartText:        item.StartText,
			EndText:          item.EndText,
		}

		tiss, err := CreateIssues(r)
		if err != nil {
			log.Fatal(err)
		}

		iss = append(iss, tiss...)
	}
}

func (c *GitHub) getMilestoneNumber(mt, owner, repository string) (string, Milestone, error) {
	m := Milestone{}

	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    id
    milestones(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        id
        number
        title
      }
    }
  }
}
`)
	req.Var("owner", owner)
	req.Var("name", repository)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResponseType
	if err := client.Run(ctx, req, &respData); err != nil {
		return "", m, err
	}

	for _, m := range respData.Repository.Milestones.Nodes {
		if m.Title == mt {
			return respData.Repository.ID, m, nil
		}
	}

	return "", m, fmt.Errorf("Error: '%s' milestone title not found", mt)
}
