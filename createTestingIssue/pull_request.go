package createTestingIssue

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
)

type Config struct {
	GitHub GitHubConfig
}

type GitHubConfig struct {
	Token      string
	Owner      string
	Repository string
	Label      string
	StartText  string
	EndText    string
}

type PullRequestResposeType struct {
	Repository Repository `json:"repository"`
}

type Repository struct {
	PullRequests PullRequestNodes `json:"pullRequests"`
	ID           string           `json:"id"`
}

type PullRequestNodes struct {
	Nodes []PullRequest `json:"nodes"`
}

type PullRequest struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	URL      string    `json:"url"`
	MergedAt time.Time `json:"mergedAt"`
}

type Request struct {
	Token      string
	Owner      string
	Repository string
	StartText  string
	EndText    string
}

func CreateIssues(r Request) ([]Issue, error) {
	var config Config

	config.GitHub.Token = r.Token
	config.GitHub.Owner = r.Owner
	config.GitHub.Repository = r.Repository
	config.GitHub.StartText = r.StartText
	config.GitHub.EndText = r.EndText

	iss := []Issue{}

	rid, prs, err := config.GitHub.getPullRequests()
	if err != nil {
		return iss, err
	}

	for _, pr := range prs {
		cir := CreateIssuesRequest{
			RepositoryID: rid,
			Title:        pr.Title,
			Body:         pr.Body,
			URL:          pr.URL,
		}
		is, err := config.GitHub.CreateIssue(cir)
		if err != nil {
			return iss, err
		}
		iss = append(iss, is)
	}

	return iss, nil
}

func (c *GitHubConfig) getPullRequests() (string, []PullRequest, error) {
	client := graphql.NewClient("https://api.github.com/graphql")

	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    id
    pullRequests(first: 100, states: [MERGED], labels: ["テストを実施"], orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        id
        title
		url
        body
		mergedAt
      }
    }
  }
}
`)

	req.Var("owner", c.Owner)
	req.Var("name", c.Repository)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData PullRequestResposeType
	err := client.Run(ctx, req, &respData)

	if err != nil {
		return "", nil, err
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", nil, err
	}
	nowJST := time.Now().In(loc)
	items := []PullRequest{}

	for _, pr := range respData.Repository.PullRequests.Nodes {
		duration := nowJST.Sub(pr.MergedAt.In(loc))
		if int(duration.Hours()) <= 24 {
			// 24時間以内にマージされたものだけをissueを作成する
			items = append(items, pr)
		}
	}

	return respData.Repository.ID, items, nil
}
