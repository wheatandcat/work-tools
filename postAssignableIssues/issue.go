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
	Categories []string
}

type IssuesNodes struct {
	Nodes []Issue `json:"nodes"`
}

type Issue struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Labels    LabelNodes `json:"labels"`
	Milestone Milestone  `json:"milestone"`
}

type LabelNodes struct {
	Nodes []Label `json:"nodes"`
}

type Label struct {
	Name string `json:"name"`
}

type Milestone struct {
	Title string `json:"title"`
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
	Categories []string
}

type TempleIssueData struct {
	Category string
	Issue    Issue
}

type TempleData struct {
	Repository string
	Issue      []TempleIssueData
}

func GetIssueText(r Request) (string, error) {
	gc := GitHubConfig(r)

	iss, err := gc.GetIssue()
	if err != nil {
		return "", err
	}

	t := TempleData{
		Repository: r.Repository,
	}

	for _, is := range iss {
		for _, c := range r.Categories {
			if is.isTempleIssueLabel(c) {
				tis := TempleIssueData{
					Category: c,
					Issue:    is,
				}
				t.Issue = append(t.Issue, tis)
			}
		}
	}

	ids := t.GetIssueIds()

	for _, is := range iss {
		if !include(ids, is.ID) {
			tis := TempleIssueData{
				Category: "開発issue",
				Issue:    is,
			}
			t.Issue = append(t.Issue, tis)
		}
	}

	text, err := t.ToBody()
	if err != nil {
		return "", err
	}

	return text, nil
}

func (t *TempleData) GetIssueIds() []string {
	ids := []string{}
	for _, is := range t.Issue {
		ids = append(ids, is.Issue.ID)
	}

	return ids
}

func (t *TempleData) GetCategories() []string {
	categories := []string{}
	for _, is := range t.Issue {
		if !include(categories, is.Category) {
			categories = append(categories, is.Category)
		}
	}

	return categories
}

func (is *Issue) isTempleIssueLabel(Label string) bool {
	for _, la := range is.Labels.Nodes {
		if la.Name == Label {
			return true
		}
	}

	return false
}

func (c *GitHubConfig) GetIssue() ([]Issue, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    issues(first: 100, states: [OPEN], labels: ["アサイン募集中"], orderBy: {field: CREATED_AT, direction: DESC}, filterBy: {assignee: null}) {
      nodes {
        id
        title
        url
        labels(first: 5, orderBy: {field: CREATED_AT, direction: DESC}) {
          nodes {
            name
          }
        }
       	milestone {
          id
          title
        }
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
	text := "### ■ " + t.Repository + " "

	if len(t.Issue) == 0 {
		return t.ToBodyNoData()
	}

	categories := t.GetCategories()

	for _, category := range categories {
		tc, err := t.ToCategoryBody(category)
		if err != nil {
			return "", err
		}

		text += tc

	}

	return text, nil
}

func (t *TempleData) ToCategoryBody(category string) (string, error) {

	type templeIssue struct {
		Category string
		Issues   []Issue
	}

	ti := templeIssue{
		Category: category,
	}

	for _, is := range t.Issue {
		if is.Category == category {
			ti.Issues = append(ti.Issues, is.Issue)
		}
	}

	if len(ti.Issues) == 0 {
		return "", nil
	}

	text := `
#### {{ .Category }}
{{range $index, $is := .Issues}}
 - {{ if $is.Milestone.Title }}【{{$is.Milestone.Title}}】{{ end }}[{{ $is.Title }}]({{ $is.URL }}){{end}}
`

	tpl, err := template.New("").Parse(text)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	if err := tpl.Execute(&doc, ti); err != nil {
		return "", err
	}

	return doc.String(), nil
}

func (t *TempleData) ToBodyNoData() (string, error) {
	text := `
募集中のissueがありません
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

func include(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
