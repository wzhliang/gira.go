package context

import (
	"github.com/davecgh/go-spew/spew"
)

// Context ...
type Context struct {
	Issue struct {
		ID          string
		FixVersions []string
		Status      string
		Summary     string
		URL         string
		Owner       string
		Components  []string
		Project     string
		HasChild    bool
	}
	PR struct {
		ID           string // for compatibility
		Owners       []string
		URL          string
		TargetBranch string
		Title        string
	}
	Sandbox       string
	WorkingDir    string
	CurrentBranch string
	Repo          struct {
		Owner string
		Name  string
	}
}

func (c *Context) Show() {
	spew.Dump(c)
}
