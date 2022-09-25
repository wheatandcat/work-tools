package main

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
)

type Config struct {
	GitHub GitHubConfig
}

type GitHubConfig struct {
	Token            string
	Owner            string
	Repository       string
	Milestone        Milestone
	TestMilestone    Milestone
	TestRepositoryID string
	Label            string
	StartText        string
	EndText          string
}

type PullRequestResponseType struct {
	Repository Repository `json:"repository"`
}

type Repository struct {
	PullRequests PullRequestNodes `json:"pullRequests"`
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	URL          string           `json:"url"`
	Milestone    Milestone        `json:"milestone"`
	Milestones   MilestoneNodes   `json:"milestones"`
}

type MilestoneNodes struct {
	Nodes []Milestone `json:"nodes"`
}

type Milestone struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Number       int              `json:"number"`
	URL          string           `json:"url"`
	PullRequests PullRequestNodes `json:"pullRequests"`
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
	Token            string
	Owner            string
	Label            string
	Milestone        Milestone
	TestMilestone    Milestone
	TestRepositoryID string
	Repository       string
	StartText        string
	EndText          string
}

func CreateIssues(r Request) ([]Issue, error) {
	var config Config

	config.GitHub.Token = r.Token
	config.GitHub.Owner = r.Owner
	config.GitHub.Repository = r.Repository
	config.GitHub.Milestone = r.Milestone
	config.GitHub.TestMilestone = r.TestMilestone
	config.GitHub.TestRepositoryID = r.TestRepositoryID
	config.GitHub.StartText = r.StartText
	config.GitHub.EndText = r.EndText
	config.GitHub.Label = r.Label

	iss := []Issue{}

	_, prs, err := config.GitHub.getPullRequests()
	if err != nil {
		return iss, err
	}

	for _, pr := range prs {
		cir := CreateIssuesRequest{
			RepositoryID: config.GitHub.TestRepositoryID,
			Title:        pr.Title,
			Body:         pr.Body,

			URL: pr.URL,
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
query Repository($owner: String!, $name: String!, $labels: [String!], $milestoneNumber: Int!) {
  repository(owner: $owner, name: $name) {
    id
    name
    url
    milestone(number: $milestoneNumber) {
      id
      title
      url
      pullRequests(first: 100, labels: $labels) {
        nodes {
          id
          title
          url
          body
        }
      }
    }
  }
}

`)

	req.Var("owner", c.Owner)
	req.Var("name", c.Repository)
	req.Var("labels", [1]string{c.Label})
	req.Var("milestoneNumber", c.Milestone.Number)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData PullRequestResponseType
	err := client.Run(ctx, req, &respData)

	if err != nil {
		return "", nil, err
	}

	return respData.Repository.ID, respData.Repository.Milestone.PullRequests.Nodes, nil
}
