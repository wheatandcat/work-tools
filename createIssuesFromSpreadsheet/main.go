package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/machinebox/graphql"
)

// Config 設定ファイルのタイプ
type Config struct {
	GitHub GitHubConfig
}

// GitHubConfig GitHub設定ファイルのタイプ
type GitHubConfig struct {
	Token      string `toml:"token"`
	Owner      string `toml:"owner"`
	Repository string `toml:"repository"`
	URL        string `toml:"url"`
}

// ResponseData GASのResponse
type ResponseData struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Milestone    string   `json:"milestone"`
	Label        string   `json:"label"`
	Repositories []string `json:"repositories"`
}

// CreateIssueResponse CreateIssueResponseのタイプ
type CreateIssueResponse struct {
	CreateIssue CreateIssueType `json:"createIssue"`
}

// CreateIssueType mutationのタイプ
type CreateIssueType struct {
	Issue Issue `json:"issue"`
}

// Issue issueのタイプ
type Issue struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ResposeType Graphqlのタイプ
type ResposeType struct {
	Repository Repository `json:"repository"`
}

// Repository Repositoryのタイプ
type Repository struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Milestones MilestoneNodes `json:"milestones"`
	Labels     LabelNodes     `json:"labels"`
}

// MilestoneNodes MilestoneNodesのタイプ
type MilestoneNodes struct {
	Nodes []Milestone `json:"nodes"`
}

// Milestone Milestoneのタイプ
type Milestone struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// LabelNodes ラベルリスト
type LabelNodes struct {
	Nodes []Label `json:"nodes"`
}

// Label ラベル
type Label struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RepositoryInfo リポジトリの情報
type RepositoryInfo struct {
	RepositoryID string
	MilestoneID  string
	Labels       []Label
}

// CreateIssue issue作成
type CreateIssue struct {
	RepositoryID string
	Title        string
	Body         string
	MilestoneID  string
	LabelIds     []string
}

// Template Templateのタイプ
type Template struct {
	ID    string
	Title string
	Body  string
	Code  string
}

// LogFile LogFileのタイプ
type LogFile struct {
	ID             string
	RepositoryName string
	URL            string
}

func include(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func main() {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	m := "v2.0.0"

	res, err := http.Get(config.GitHub.URL + "?milestone=" + m)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("StatusCode=%d", res.StatusCode)
		log.Fatal("OUT")
	}

	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		log.Fatal(error)
	}

	var data []ResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	ri, err := config.GitHub.getRepositoryInfo(m)
	if err != nil {
		log.Fatal(err)
	}

	var logFiles []LogFile
	raw, err := ioutil.ReadFile("./log.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(raw, &logFiles)

	var logIDList []string
	for _, logFile := range logFiles {
		if logFile.RepositoryName == config.GitHub.Repository {
			logIDList = append(logIDList, logFile.ID)
		}
	}

	for _, item := range data {
		if !include(item.Repositories, config.GitHub.Repository) {
			continue
		}

		if include(logIDList, strconv.Itoa(item.ID)) {
			// 既に作成済み
			fmt.Println("ID: " + strconv.Itoa(item.ID) + "は作成済みです(" + config.GitHub.Repository + ")")
			continue
		}

		lid := ""
		for _, label := range ri.Labels {
			if label.Name == item.Label {
				lid = label.ID
			}
		}

		title := strconv.Itoa(item.ID) + "_" + item.Title

		tmp := Template{
			ID:    strconv.Itoa(item.ID),
			Title: item.Title,
			Body:  item.Body,
			Code:  "```",
		}

		body, err := tmp.ToBody()
		if err != nil {
			log.Fatal(err)
		}

		ci := CreateIssue{
			RepositoryID: ri.RepositoryID,
			Title:        title,
			Body:         body,
			MilestoneID:  ri.MilestoneID,
			LabelIds:     []string{lid},
		}

		cires, err := config.GitHub.createIssue(ci)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("ID: " + strconv.Itoa(item.ID) + "は作成しました(" + config.GitHub.Repository + ")")

		logFile := LogFile{
			ID:             strconv.Itoa(item.ID),
			RepositoryName: config.GitHub.Repository,
			URL:            cires.CreateIssue.Issue.URL,
		}
		logFiles = append(logFiles, logFile)
	}

	file, _ := json.MarshalIndent(logFiles, "", " ")
	_ = ioutil.WriteFile("log.json", file, 0644)

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
    labels(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
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
		return ri, fmt.Errorf("Error: '%s' milestone title not found", mt)
	}

	for _, l := range respData.Repository.Labels.Nodes {
		ri.Labels = append(ri.Labels, l)
	}

	ri.MilestoneID = mid
	ri.RepositoryID = rid

	return ri, nil
}

func (c *GitHubConfig) createIssue(ci CreateIssue) (CreateIssueResponse, error) {
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

// ToBody 本文作成
func (t *Template) ToBody() (string, error) {
	text := `
## 概要

No.{{ .ID }} {{ .Title }}

{{ .Code }}
{{ .Body }}
{{ .Code }}

## 対応内容

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
