package main

import (
	"bytes"
	"context"
	"strings"
	"text/template"

	"github.com/machinebox/graphql"
)

type CreateIssuesRequest struct {
	RepositoryID string
	Title        string
	Body         string
	URL          string
}

type TempleData struct {
	Title string
	Body  string
	URL   string
}

type ResponseType struct {
	Repository Repository `json:"repository"`
}

type CreateIssueResponse struct {
	CreateIssue CreateIssueType `json:"createIssue"`
}

type CreateIssueType struct {
	Issue Issue `json:"issue"`
}

type Issue struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (c *GitHubConfig) CreateIssue(r CreateIssuesRequest) (Issue, error) {
	start := strings.Index(r.Body, c.StartText)
	if start == -1 {
		start = 0
	} else {
		start += len(c.StartText)
	}
	end := strings.Index(r.Body, c.EndText)
	if end == -1 {
		end = len(r.Body)
	}
	if end < start {
		end = len(r.Body)
	}

	body := r.Body[start:end]

	text := `
## 対象issue
 - {{ .URL }}

## テスト内容
{{ .Body }}
`
	tmp, err := template.New("").Parse(text)
	if err != nil {
		return Issue{}, err
	}

	td := TempleData{
		Title: r.Title,
		Body:  body,
		URL:   r.URL,
	}

	var doc bytes.Buffer
	if err := tmp.Execute(&doc, td); err != nil {
		return Issue{}, err
	}

	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
mutation CreateIssue($repositoryId: ID!, $title: String!, $body: String,$milestoneId: ID!) {
  createIssue(input: {repositoryId: $repositoryId, title: $title, body: $body, milestoneId: $milestoneId}) {
    issue {
      id
      title
      url
    }
  }
}

`)

	req.Var("repositoryId", r.RepositoryID)
	req.Var("title", r.Title)
	req.Var("body", doc.String())
	req.Var("milestoneId", c.TestMilestone.ID)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData CreateIssueResponse
	if err := client.Run(ctx, req, &respData); err != nil {
		return Issue{}, err
	}

	return respData.CreateIssue.Issue, nil
}
