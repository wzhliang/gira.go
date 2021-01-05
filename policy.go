package main

import (
	"fmt"
	"time"

	"github.com/wzhliang/gira/pkg/context"
	"github.com/wzhliang/gira/pkg/git"
)

// Enforcer ...
type Enforcer interface {
	Check(c *context.Context) bool
	Message(c *context.Context) string
	ErrMessage(c *context.Context) string
}

// Policy ...
type Policy struct {
	pols []Enforcer
}

// Add ...
func (p Policy) Add(e Enforcer) Policy {
	p.pols = append(p.pols, e)
	return p
}

// Check ...
func (p Policy) Check(ctx *context.Context) bool {
	ret := true
	for _, e := range p.pols {
		if e.Message(ctx) != "" {
			fmt.Printf("%s ...", e.Message(ctx))
		}
		if !e.Check(ctx) {
			fmt.Printf("    ❌\n")
			fmt.Printf("→    %s\n", e.ErrMessage(ctx))
			ret = false
			break
		} else {
			fmt.Printf("    ✔️\n")
		}
	}
	return ret
}

// IssueStatusChecker ...
type IssueStatusChecker struct{}

// Check ...
func (ic IssueStatusChecker) Check(ctx *context.Context) bool {
	return _jra.IsOpen(ctx)
}

// Message ...
func (ic IssueStatusChecker) Message(ctx *context.Context) string {
	return "Issue can be started"
}

// ErrMessage ...
func (ic IssueStatusChecker) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Issue %s in wrong state (%s) and cannot be started.", ctx.Issue.ID, ctx.Issue.Status)
}

// IssueTypeChecker ...
type IssueTypeChecker struct{}

// Check ...
func (ic IssueTypeChecker) Check(ctx *context.Context) bool {
	return ctx.Issue.HasChild == false
}

// Message ...
func (ic IssueTypeChecker) Message(ctx *context.Context) string {
	return "Issue has no subtask"
}

// ErrMessage ...
func (ic IssueTypeChecker) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Issue %s is either an Epic or has subtasks. Cannot be started.", ctx.Issue.ID)
}

// IssueAssigned ...
type IssueAssigned struct{}

// Check ...
func (ic IssueAssigned) Check(ctx *context.Context) bool {
	return ctx.Issue.Owner != ""
}

// Message ...
func (ic IssueAssigned) Message(ctx *context.Context) string {
	return "Issue has been assigned"
}

// ErrMessage ...
func (ic IssueAssigned) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Issue not assigned yet.")
}

// IssuePlanned ...
type IssuePlanned struct{}

// Check ...
func (ic IssuePlanned) Check(ctx *context.Context) bool {
	return len(ctx.Issue.FixVersions) > 0
}

// Message ...
func (ic IssuePlanned) Message(ctx *context.Context) string {
	return "Issue has been planned"
}

// ErrMessage ...
func (ic IssuePlanned) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Issue not planned yet (empty Fix Version/s)")
}

// WaitForBranch ...
type WaitForBranch struct{}

// Check ...
func (w WaitForBranch) Check(ctx *context.Context) bool {
	ok := false
	for i := 0; i < 3; i = i + 1 {
		err := _gte.GetBranch(ctx)
		if err == nil {
			ok = true
			break
		}
		time.Sleep(5 * time.Second)
	}
	return ok
}

// Message ...
func (w WaitForBranch) Message(ctx *context.Context) string {
	return fmt.Sprintf("Waiting for remote branch %s to be available", ctx.Issue.ID)
}

// ErrMessage ...
func (w WaitForBranch) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Timeout waiting for remote branch %s.", ctx.Issue.ID)
}

// JiraChecker ...
type JiraChecker struct {
}

// Check ...
func (jc JiraChecker) Check(ctx *context.Context) bool {
	return !_jra.IsDone(ctx)
}

// Message ...
func (jc JiraChecker) Message(ctx *context.Context) string {
	return "Issue is correct state"
}

// ErrMessage ...
func (jc JiraChecker) ErrMessage(ctx *context.Context) string {
	return "Issue is already finished."
}

// IssueChecker ...
type IssueChecker struct {
}

// Check ...
func (jc IssueChecker) Check(ctx *context.Context) bool {
	return _ctx.Issue.ID != ""
}

// Message ...
func (jc IssueChecker) Message(ctx *context.Context) string {
	return "Issue ID is valid"
}

// ErrMessage ...
func (jc IssueChecker) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Invalid issue ID: %s", ctx.Issue.ID)
}

// IssueComponentChecker ...
type IssueComponentChecker struct {
}

// Check ...
func (jc IssueComponentChecker) Check(ctx *context.Context) bool {
	return len(ctx.Issue.Components) > 0
}

// Message ...
func (jc IssueComponentChecker) Message(ctx *context.Context) string {
	return "Issue has valid components"
}

// ErrMessage ...
func (jc IssueComponentChecker) ErrMessage(ctx *context.Context) string {
	return "Component for issue has not yet been identified."
}

// IssueMergable ...
type IssueMergable struct {
	msg string
}

// Check ...
func (jc *IssueMergable) Check(ctx *context.Context) bool {
	jc.msg = "Unknown error."
	if _ctx.Issue.ID == "" {
		jc.msg = "Empty issue ID."
		return false
	}
	if _jra.IsDone(ctx) {
		jc.msg = "Issue already finished."
		return false
	}
	if len(ctx.Issue.FixVersions) == 0 {
		jc.msg = "Issue not planned (empty Fix version/s)."
		return false
	}
	// Check epic
	return true
}

// Message ...
func (jc *IssueMergable) Message(ctx *context.Context) string {
	return "Issue is mergable"
}

// ErrMessage ...
func (jc *IssueMergable) ErrMessage(ctx *context.Context) string {
	return jc.msg
}

// PROwnerChecker ...
type PROwnerChecker struct {
}

// Check ...
func (prc PROwnerChecker) Check(ctx *context.Context) bool {
	return len(_ctx.PR.Owners) >= 2
}

// Message ...
func (prc PROwnerChecker) Message(ctx *context.Context) string {
	return "PR is assigned to both tester and reivewer"
}

// ErrMessage ...
func (prc PROwnerChecker) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("PR assgigned to %d people instead of 2.", len(_ctx.PR.Owners))
}

// PRTitleChecker ...
type PRTitleChecker struct {
}

// Check ...
func (prc PRTitleChecker) Check(ctx *context.Context) bool {
	return _jra.ValidIssueID(ctx)
}

// Message ...
func (prc PRTitleChecker) Message(ctx *context.Context) string {
	return "PR title is valid"
}

// ErrMessage ...
func (prc PRTitleChecker) ErrMessage(ctx *context.Context) string {
	return fmt.Sprintf("Invalid PR title (%s). Should start with JIRA issue ID", ctx.PR.Title)
}

// CleanTreeChecker ...
type CleanTreeChecker struct{}

// Check ...
func (cc CleanTreeChecker) Check(ctx *context.Context) bool {
	return git.StashableChanges() == false
}

// Message ...
func (cc CleanTreeChecker) Message(ctx *context.Context) string {
	return "Local tree is clean"
}

// ErrMessage ...
func (cc CleanTreeChecker) ErrMessage(ctx *context.Context) string {
	return "Local tree is not clean. git status to see"
}

// PRBranchChecker ...
type PRBranchChecker struct {
	msg string
}

// Check ...
func (pc *PRBranchChecker) Check(ctx *context.Context) bool {
	br := git.CurrentBranch()
	if br == "master" {
		pc.msg = "You have to be on PR branch to create a PR"
		return false
	}
	rebase, err := git.NeedsRebase(br, "master")
	if err != nil || rebase {
		pc.msg = "Looks like your branch requires rebasing"
		return false
	}

	return true
}

// Message ...
func (pc *PRBranchChecker) Message(ctx *context.Context) string {
	return "PR branch is valid and up-to-date"
}

// ErrMessage ...
func (pc *PRBranchChecker) ErrMessage(ctx *context.Context) string {
	return pc.msg
}
