package jira

import (
	"fmt"
	"regexp"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/rs/zerolog/log"
	"github.com/wzhliang/gira/pkg/context"
)

// Jira ...
type Jira struct {
	User   string
	Passwd string
	Root   string
	cli    *jira.Client
}

// New ...
func New(url, user, passwd string) *Jira {
	log.Info().Msgf("JIRA::New(%s, %s, %s)\n", url, user, passwd)
	tp := jira.BasicAuthTransport{
		Username: user,
		Password: passwd,
	}

	cli, err := jira.NewClient(tp.Client(), url)
	if err != nil {
		return nil
	}

	return &Jira{
		User:   user,
		Passwd: passwd,
		Root:   url,
		cli:    cli,
	}
}

// GetIssue initiliases everything about a JIRA issue
func (j *Jira) GetIssue(ctx *context.Context) error {
	log.Info().Msgf("GetIssue %s", ctx.Issue.ID)
	iss, _, err := j.cli.Issue.Get(ctx.Issue.ID, nil)
	if err != nil {
		return err
	}

	ctx.Issue.Summary = iss.Fields.Summary
	// we only care about the category
	ctx.Issue.Status = iss.Fields.Status.StatusCategory.Name
	// TODO: take care of "//"
	ctx.Issue.URL = fmt.Sprintf("%s/browse/%s", j.Root, ctx.Issue.ID)
	if iss.Fields.Assignee != nil {
		ctx.Issue.Owner = iss.Fields.Assignee.Name
	}
	for _, ver := range iss.Fields.FixVersions {
		ctx.Issue.FixVersions = append(ctx.Issue.FixVersions, ver.Name)
	}
	for _, cp := range iss.Fields.Components {
		ctx.Issue.Components = append(ctx.Issue.Components, cp.Name)
	}
	ctx.Issue.Project = strings.Split(ctx.Issue.ID, "-")[0]
	ctx.Issue.HasChild = (iss.Fields.Type.Name == "Epic") || (len(iss.Fields.Subtasks) > 0)
	return nil
}

// IssueStatus ...
func (j *Jira) IssueStatus(ctx *context.Context) string {
	return ctx.Issue.Status
}

// IsDone ...
func (j *Jira) IsDone(ctx *context.Context) bool {
	return ctx.Issue.Status == "Done"
}

// IsOpen ...
func (j *Jira) IsOpen(ctx *context.Context) bool {
	return ctx.Issue.Status == "To Do"
}

// ValidIssueID ...
func (j *Jira) ValidIssueID(ctx *context.Context) bool {
	matched, err := regexp.MatchString(`[a-zA-Z]+-\d+`, ctx.Issue.ID)
	if err != nil {
		return false
	}

	return matched
}

// UpdateIssue ..
func (j *Jira) UpdateIssue(ctx *context.Context, comment, transition string) error {
	log.Info().Msgf("UpdateIssue(%s, %s, %s)", ctx.Issue.ID, comment, transition)
	_, _, err := j.cli.Issue.AddComment(ctx.Issue.ID, &jira.Comment{Body: comment})
	if err != nil {
		return err
	}
	_, err = j.cli.Issue.DoTransition(ctx.Issue.ID, transition)
	if err != nil {
		return err
	}
	return nil
}

// List ...
func (j *Jira) List(ctx *context.Context) ([]string, error) {
	log.Info().Msg("JIRA:List...")
	var ret []string
	var jql string
	if ctx.Issue.Project != "" {
		jql = fmt.Sprintf(`project=%s and assignee=%s and statusCategory="To Do"`, ctx.Issue.Project, j.User)
	} else {
		jql = fmt.Sprintf(`assignee=%s and statusCategory="To Do"`, j.User)
	}
	log.Debug().Msgf("JQL: %s", jql)
	list, _, err := j.cli.Issue.Search(jql, &jira.SearchOptions{})
	if err != nil {
		return nil, err
	}

	for _, iss := range list {
		ret = append(ret, fmt.Sprintf("%s - %s", iss.Key, iss.Fields.Summary))
	}

	log.Debug().Msgf("ret: %v", ret)
	return ret, nil
}
