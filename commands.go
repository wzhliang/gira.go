package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/wzhliang/gira/pkg/config"
	"github.com/wzhliang/gira/pkg/context"
	"github.com/wzhliang/gira/pkg/gitee"
	"github.com/wzhliang/gira/pkg/jira"
)

// CmdContext ...
type CmdContext struct {
	jra  *jira.Jira
	gte  *gitee.Gitee
	ctx  *context.Context
	conf *config.Config
}

var cli struct {
	Merge  MergeCmd  `cmd help:"Merge PR and resolve relating issue."`
	Start  StartCmd  `cmd help:"Start working on an JIRA issue."`
	Finish FinishCmd `cmd help:"Finis JIRA and create PR."`
	Test   TestCmd   `cmd help:"Handy Tests"`
}

func kongMain(conf *config.Config, j *jira.Jira, g *gitee.Gitee, c *context.Context) {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&CmdContext{
		jra:  j,
		gte:  g,
		ctx:  c,
		conf: conf,
	})
	if err != nil {
		fmt.Printf("%v", err)
	}
}

// Info ...
func Info(msg string) {
	fmt.Printf("►  %s\n", msg)
}
