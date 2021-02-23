package main

// For accessing globals....

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/wzhliang/gira/pkg/config"
	"github.com/wzhliang/gira/pkg/git"
)

// FinishCmd ...
type FinishCmd struct {
	IssueID string `arg optional name:"issue"`
}

func getTargetBranch(versions []string) (string, error) {
	// FIXME: properly implement when CP is supported
	return "master", nil
}

// Run ...
func (sc *FinishCmd) Run(cmd *CmdContext) error {
	c := cmd.ctx

	c.Issue.ID = sc.IssueID
	if sc.IssueID == "" {
		c.Issue.ID = c.CurrentBranch
	}

	fmt.Printf("Creating PR for %s\n", c.Issue.ID)
	err := cmd.jra.GetIssue(c)
	if err != nil {
		fmt.Printf("Unable to get JIRA issue: %s\n", c.Issue.ID)
		fmt.Println(err)
	}

	cont := Policy{}.
		Add(Enforcer(CleanTreeChecker{})).
		Add(Enforcer(&PRBranchChecker{})).
		Check(c)
	if !cont {
		return errors.New("----")
	}

	fmt.Println("")
	Info("git pushing...")
	// TODO: check if rebasing is required
	err = git.Push()
	if err != nil {
		fmt.Printf("Unable to push")
	}
	c.PR.TargetBranch, err = getTargetBranch(c.Issue.FixVersions)
	if err != nil {
		fmt.Printf("Unable to figure out target branch for PR.")
		fmt.Printf("Fix version/s: %v", c.Issue.FixVersions)
		return err
	}

	Info("Creating PR...")
	err = cmd.gte.CreatePR(c)
	if err != nil {
		fmt.Printf("Uanble to create PR.")
		return err
	}
	Info("Updating JIRA issue...")
	conf := reflect.Indirect(reflect.ValueOf(cmd.conf)).FieldByName(c.Issue.Project).Interface()
	err = cmd.jra.UpdateIssue(c,
		fmt.Sprintf("PR created: %s", c.PR.URL),
		fmt.Sprintf("%d", conf.(config.Status).Ready_for_test))
	if err != nil {
		fmt.Printf("Uanble to update issue.")
		return err
	}

	return nil
}
