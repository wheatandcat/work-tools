package postAssignableIssues

import (
	"context"

	"github.com/machinebox/graphql"
)

type RequestNote struct {
	Token   string
	ID      string
	Host    string
	Content string
}

type Note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type ResposeNoteType struct {
	Note Note `json:"note"`
}

type UpdateNoteContentInput struct {
	ID          string `json:"id"`
	NewContent  string `json:"newContent"`
	BaseContent string `json:"baseContent"`
}

func PostNote(r RequestNote) error {
	n, err := r.GetNote()
	if err != nil {
		return err
	}
	input := UpdateNoteContentInput{
		ID:          r.ID,
		NewContent:  r.Content,
		BaseContent: n.Content,
	}
	if err := r.UpdateNote(input); err != nil {
		return err
	}

	return nil
}

func (c *RequestNote) GetNote() (Note, error) {
	client := graphql.NewClient("https://" + c.Host + ".kibe.la/api/v1")
	req := graphql.NewRequest(`
query Note($id: ID!) {
  note(id: $id) {
    id
    title
    url
    content
  }
}
`)
	req.Var("id", c.ID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResposeNoteType
	if err := client.Run(ctx, req, &respData); err != nil {
		return Note{}, err
	}

	return respData.Note, nil
}

func (c *RequestNote) UpdateNote(input UpdateNoteContentInput) error {
	client := graphql.NewClient("https://" + c.Host + ".kibe.la/api/v1")
	req := graphql.NewRequest(`
mutation UpdateNoteContent($input: UpdateNoteContentInput!) {
  updateNoteContent(input: $input) {
    clientMutationId
    note {
      id
      content
    }
  }
}
`)
	req.Var("input", input)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResposeNoteType
	if err := client.Run(ctx, req, &respData); err != nil {
		return err
	}

	return nil
}
