package createissue

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"

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

type ResponseData struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Image        string   `json:"image"`
	Env          string   `json:"env"`
	Priority     string   `json:"priority"`
	Repositories []string `json:"repositories"`
}

type CreateIssueResponse struct {
	CreateIssue CreateIssueType `json:"createIssue"`
	ID          int             `json:"id"`
}

type CreateIssueType struct {
	Issue Issue `json:"issue"`
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
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Milestones MilestoneNodes `json:"milestones"`
	Labels     LabelNodes     `json:"labels"`
}

type MilestoneNodes struct {
	Nodes []Milestone `json:"nodes"`
}

type Milestone struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type LabelNodes struct {
	Nodes []Label `json:"nodes"`
}

type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RepositoryInfo struct {
	RepositoryID string
	MilestoneID  string
	Labels       []Label
}

type CreateIssueRequest struct {
	RepositoryID string
	Title        string
	Body         string
	MilestoneID  string
	LabelIds     []string
}

type Template struct {
	ID    string
	Row   string
	Title string
	Body  string
	Code  string
	Image string
	Env   string
}

type Request struct {
	Token      string
	Owner      string
	Version    string
	Repository string
	ID         int
	Priority   string
	Title      string
	Body       string
	Env        string
	Image      string
}

func Create(r Request) (CreateIssueResponse, error) {

	var cir CreateIssueResponse

	gc := GitHubConfig{
		Token:      r.Token,
		Owner:      r.Owner,
		Repository: r.Repository,
	}

	ri, err := gc.getRepositoryInfo(r.Version)
	if err != nil {
		return cir, err
	}

	bugLabel := ""
	testLabel := ""
	for _, label := range ri.Labels {
		if label.Name == "テスト差し戻し" {
			testLabel = label.ID
		} else if label.Name == "バグ" {
			bugLabel = label.ID
		}
	}

	title := strconv.Itoa(r.ID) + "_" + r.Title

	tmp := Template{
		ID:    strconv.Itoa(r.ID),
		Row:   strconv.Itoa(r.ID + 1),
		Body:  r.Body,
		Image: r.Image,
		Env:   r.Env,
		Code:  "```",
	}
	body, err := tmp.ToBody()
	if err != nil {
		return cir, err
	}

	ci := CreateIssueRequest{
		RepositoryID: ri.RepositoryID,
		Title:        title,
		Body:         body,
		MilestoneID:  ri.MilestoneID,
		LabelIds:     []string{bugLabel, testLabel},
	}

	if r.Priority == "高" {
		cir, err = gc.createIssue(ci)
		if err != nil {
			return cir, err
		}
	} else {
		cir, err = gc.createIssue2(ci)
		if err != nil {
			return cir, err
		}
	}

	cir.ID = r.ID

	return cir, nil
}

func (c *GitHubConfig) getRepositoryInfo(mt string) (RepositoryInfo, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    id
    milestones(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        id
        title
      }
    }
    labels(first: 30, orderBy: {field: CREATED_AT, direction: ASC}) {
      nodes {
        id
        name
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

	ri := RepositoryInfo{}

	var respData ResposeType
	if err := client.Run(ctx, req, &respData); err != nil {
		return ri, err
	}

	rid := respData.Repository.ID
	mid := ""

	for _, m := range respData.Repository.Milestones.Nodes {
		if m.Title == mt {
			mid = m.ID
		}
	}

	if mid == "" {
		return ri, fmt.Errorf("'%s' milestone title not found", mt)
	}

	ri.Labels = append(ri.Labels, respData.Repository.Labels.Nodes...)
	ri.MilestoneID = mid
	ri.RepositoryID = rid

	return ri, nil
}

func (c *GitHubConfig) createIssue(ci CreateIssueRequest) (CreateIssueResponse, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
mutation CreateIssue(
	$repositoryId: ID!
	$title: String!
	$body: String
	$milestoneId: ID!
	$labelIds: [ID!]
  ) {
  createIssue(input: {repositoryId: $repositoryId, title:$title, body: $body, milestoneId: $milestoneId, labelIds: $labelIds}) {
    issue {
	  id
	  title
	  url
    }
  }
}
`)

	req.Var("repositoryId", ci.RepositoryID)
	req.Var("title", ci.Title)
	req.Var("body", ci.Body)
	req.Var("milestoneId", ci.MilestoneID)
	req.Var("labelIds", ci.LabelIds)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData CreateIssueResponse
	if err := client.Run(ctx, req, &respData); err != nil {
		return respData, err
	}

	return respData, nil
}

func (c *GitHubConfig) createIssue2(ci CreateIssueRequest) (CreateIssueResponse, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
mutation CreateIssue(
	$repositoryId: ID!
	$title: String!
	$body: String
	$labelIds: [ID!]
  ) {
  createIssue(input: {repositoryId: $repositoryId, title:$title, body: $body, labelIds: $labelIds}) {
    issue {
	  id
	  url
    }
  }
}
`)

	req.Var("repositoryId", ci.RepositoryID)
	req.Var("title", ci.Title)
	req.Var("body", ci.Body)
	req.Var("labelIds", ci.LabelIds)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData CreateIssueResponse
	if err := client.Run(ctx, req, &respData); err != nil {
		return respData, err
	}

	return respData, nil
}

func (t *Template) ToBody() (string, error) {
	text := `
## 概要

No.{{ .ID }} {{ .Title }}

{{ .Code }}
{{ .Body }}
{{ .Code }}

■ 参考画像
{{ .Image }}

■ 発生したOS・ブラウザ
{{ .Env }}

## 対応内容

 - 上記を修正する

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
