package git

import (
	"fmt"
	"os/exec"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	parser "github.com/whilp/git-urls"
)

// Partially borrowed from https://raw.githubusercontent.com/chriswalz/bit/master/cmd/git.go

func GetRemote(name string) string {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return ""
	}

	r, err := repo.Remotes()
	if err != nil {
		return ""
	}

	for _, v := range r {
		if v.Config().Name == name {
			return v.Config().URLs[0]
		}
	}

	return ""
}

// Info ...
func Info(origin string) (owner, repo string, err error) {
	url, err := parser.Parse(origin)
	if err != nil {
		return "", "", err
	}
	// TODO: seems like an lack of feature from upstream
	// https://git-scm.com/book/en/v2/Git-on-the-Server-The-Protocols
	//    scp style: [user@]server:project.git
	tmp := strings.Split(url.Path, "/")
	owner = tmp[0]
	repo = strings.Replace(tmp[1], ".git", "", -1)
	err = nil
	return
}

// CurrentBranch ...
func CurrentBranch() string {
	msg, err := execCommand("git", "branch", "--show-current").CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	return strings.TrimSpace(string(msg))
}

// IsAheadOfCurrent ...
func IsAheadOfCurrent() bool {
	msg, err := execCommand("git", "status", "-sb").CombinedOutput()
	if err != nil {
		log.Debug().Msg(err.Error())
	}
	return strings.Contains(string(msg), "ahead")
}

// IsGitRepo ...
func IsGitRepo() bool {
	_, err := execCommand("git", "status").CombinedOutput()
	if err != nil {
		return false
	}
	return true
}

// IsBehindCurrent ...
func IsBehindCurrent() bool {
	msg, err := execCommand("git", "status", "-sb").CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	return strings.Contains(string(msg), "behind")
}

// NothingToCommit ...
func NothingToCommit() bool {
	// git diff-index HEAD --
	msg, err := execCommand("git", "diff-index", "HEAD", "--").CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	changedFiles := strings.Split(strings.TrimSpace(string(msg)), "\n")
	if len(changedFiles) == 1 && changedFiles[0] == "" {
		return true
	}
	return false
}

// IsDiverged ...
func IsDiverged() bool {
	msg, err := execCommand("git", "status").CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	return strings.Contains(string(msg), "have diverged")
}

// MostRecentCommonAncestorCommit ...
func MostRecentCommonAncestorCommit(branchA, branchB string) string {
	msg, err := execCommand("git", "merge-base", branchA, branchB).CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	return string(msg)
}

// CheckoutBranch ...
func CheckoutBranch(branch string) bool {
	msg, err := execCommand("git", "checkout", branch).CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	if strings.Contains(string(msg), "did not match any file") {
		return false
	}
	fmt.Println(string(msg))
	return true
}

// Pull ...
func Pull() error {
	_, err := execCommand("git", "pull").CombinedOutput()
	if err != nil {
		return err
	}
	log.Debug().Msg("Branch was fast-forwarded by bit.")
	return nil
}

// Push ...
func Push() error {
	_, err := execCommand("git", "push").CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// StashableChanges ...
func StashableChanges() bool {
	msg, err := execCommand("git", "status").CombinedOutput()
	if err != nil {
		log.Debug().Err(err)
	}
	return strings.Contains(string(msg), "Changes to be committed:") || strings.Contains(string(msg), "Changes not staged for commit:")
}

// GetRoot ...
func GetRoot() string {
	msg, err := execCommand("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if err != nil {
		return ""
	}
	if strings.Contains(string(msg), "fatal: not a git") {
		return ""
	}
	return strings.TrimSpace(string(msg))
}

// NeedsRebase ...
func NeedsRebase(current, upstream string) (bool, error) {
	upstreamHash, err := execCommand("git", "rev-parse", upstream).Output()
	if err != nil {
		return true, err
	}
	mergeBase, err := execCommand("git", "merge-base", current, upstream).Output()
	if err != nil {
		return true, err
	}

	return string(upstreamHash) != string(mergeBase), nil
}

func tagCurrentBranch(version string) error {
	msg, err := execCommand("git", "tag", version).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %v", string(msg), err)
	}
	return err
}

func execCommand(name string, arg ...string) *exec.Cmd {
	log.Debug().Msg(name + " " + strings.Join(arg, " "))
	c := exec.Command(name, arg...)

	/* For some reason, appending this cause failure
	if name == "git" {
		// exec commands are parsed by bit without getting printed.
		// parsing git assumes english
		c.Env = append(c.Env, "LANG=C")
	}
	*/
	return c
}
