package main

// For accessing globals....

import (
	"fmt"
)

// StartCmd ...
type TestCmd struct {
}

// Run ...
func (sc *TestCmd) Run(cmd *CmdContext) error {
	cmd.ctx.Issue.ID = "YUN-11"
	cmd.jra.GetIssue(cmd.ctx)
	list, err := cmd.jra.List(cmd.ctx)
	if err != nil {
		panic(err)
	}
	for _, iss := range list {
		fmt.Println(iss)
	}
	return nil
}
