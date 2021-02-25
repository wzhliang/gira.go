package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/wzhliang/gira/pkg/config"
	"github.com/wzhliang/gira/pkg/git"
)

// StartCmd ...
type StartCmd struct {
	IssueID string `arg optional name:issue`
}

// Run ...
func (sc *StartCmd) Run(cmd *CmdContext) error {
	if sc.IssueID == "" {
		sc.IssueID = Pick(cmd, Lister(cmd.jra), "Which issue to start?")
	}
	if sc.IssueID == "" {
		return nil
	}
	c := cmd.ctx
	c.Issue.ID = sc.IssueID
	err := cmd.jra.GetIssue(c)
	if err != nil {
		fmt.Printf("Unable to get JIRA issue information.")
		return err
	}
	// update issue happens before checking...
	// right now we're relying on webhook
	conf := reflect.Indirect(reflect.ValueOf(cmd.conf)).FieldByName(c.Issue.Project).Interface()
	cmd.jra.UpdateIssue(c, "Starting ...", fmt.Sprintf("%d", conf.(config.Status).In_progress))

	cont := Policy{}.
		Add(Enforcer(IssueStatusChecker{})).
		Add(Enforcer(IssueAssigned{})).
		Add(Enforcer(IssuePlanned{})).
		Add(Enforcer(IssueTypeChecker{})).
		Add(Enforcer(IssueComponentChecker{})).
		Add(Enforcer(WaitForBranch{})).
		Check(c)
	if !cont {
		return errors.New("")
	}

	fmt.Println()
	Info("Pulling master...")
	git.CheckoutBranch("master")
	err = git.Pull()
	if err != nil {
		fmt.Printf("git pull failed")
		fmt.Println(err)
		return err
	}
	Info("Switching to PR branch...")
	git.CheckoutBranch(c.Issue.ID)
	Info("You're all set. 请开始你的表演☀️")

	return nil
}
