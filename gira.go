package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/wzhliang/gira/pkg/config"
	"github.com/wzhliang/gira/pkg/context"
	"github.com/wzhliang/gira/pkg/git"
	"github.com/wzhliang/gira/pkg/gitee"
	"github.com/wzhliang/gira/pkg/jira"
)

var _ctx context.Context
var _jra *jira.Jira
var _gte *gitee.Gitee

func configPath() string {
	home, ok := os.LookupEnv("HOME")
	if !ok {
		panic("HOME not defined!!!")
	}

	return fmt.Sprintf("%s/.config/gira.toml", home)
}

func main() {
	if os.Getenv("GIRA_LOG") == "" {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
	conf := config.New(configPath())
	if conf == nil {
		panic("Unable to load config")
	}

	boot(conf)
	kongMain(conf, _jra, _gte, &_ctx)
}

func boot(ctx *config.Config) {
	gitMain(ctx)
	jiraMain(ctx)
}

func jiraMain(cfg *config.Config) {
	_jra = jira.New(cfg.Jira.Url, cfg.Jira.User, cfg.Jira.Passwd)
}

func gitMain(cfg *config.Config) {
	dir := git.GetRoot()
	if dir == "" {
		panic("It looks like you're not in a git repo")
	}
	os.Chdir(dir)
	_ctx.WorkingDir = dir
	orig := git.GetRemote("origin")
	if orig == "" {
		panic("Remote 'origin' not defined")
	}
	_ctx.Sandbox = orig
	_ctx.CurrentBranch = git.CurrentBranch()

	var err error
	_ctx.Repo.Owner, _ctx.Repo.Name, err = git.Info(orig)
	if err != nil {
		panic("Fail to parse remote URL")
	}
	_gte = gitee.New(_ctx.Repo.Name, _ctx.Repo.Owner, cfg.Gitee.Token)
}
