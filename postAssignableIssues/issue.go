package postAssignableIssues

import (
	"bytes"
	"context"
	"html/template"

	"github.com/machinebox/graphql"
)

type Config struct {
	GitHub GitHubConfig
}

type GitHubConfig struct {
	Token      string
	Owner      string
	Repository string
}

type IssuesNodes struct {
	Nodes []Issue `json:"nodes"`
}

type Issue struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type ResposeType struct {
	Repository Repository `json:"repository"`
}

type Repository struct {
	Issues IssuesNodes `json:"issues"`
}

type Request struct {
	Token      string
	Owner      string
	Repository string
}

type TempleData struct {
	Repository string
	Issues     []Issue
}

func GetIssueText(r Request) (string, error) {
	gc := GitHubConfig{
		Token:      r.Token,
		Owner:      r.Owner,
		Repository: r.Repository,
	}

	is, err := gc.GetIssue()
	if err != nil {
		return "", err
	}

	t := TempleData{
		Repository: r.Repository,
		Issues:     is,
	}

	text, err := t.ToBody()
	if err != nil {
		return "", err
	}

	return text, nil
}

func (c *GitHubConfig) GetIssue() ([]Issue, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    issues(first: 100, labels: ["アサイン募集中"], orderBy: {field: CREATED_AT, direction: DESC}, filterBy: {assignee: null}) {
      nodes {
        id
        title
        url
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

	var respData ResposeType
	if err := client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}

	is := make([]Issue, len(respData.Repository.Issues.Nodes))

	for i, issue := range respData.Repository.Issues.Nodes {
		is[i] = issue
	}

	return is, nil
}

func (t *TempleData) ToBody() (string, error) {
	text := `
### ■ {{ .Repository }}
{{range $index, $is := .Issues}}
 - [{{ $is.Title }}](#{{ $is.URL }}){{end}}
`

	tpl, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	if err := tpl.Execute(&doc, t); err != nil {
		return "", err
	}

	return doc.String(), nil
}
