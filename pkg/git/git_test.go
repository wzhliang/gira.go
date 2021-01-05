package git

import (
	"os"
	"strings"
	"testing"
)

func _TestPull(t *testing.T) {
	err := Pull()
	if err != nil {
		t.Log(err)
		t.Errorf("should have succeeded.")
	}
}

func TestChdir(t *testing.T) {
	path := "/Users/wliang/github/gira"
	err := os.Chdir(path)
	if err != nil {
		t.Errorf("chdir: %v", err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Getwd: %v", err)
	}
	t.Logf("pwd: %s\n", pwd)
	if path != pwd {
		t.Errorf("Getwd: %s", pwd)
	}
}

func TestRoot(t *testing.T) {
	root := GetRoot()
	if root == "" {
		t.Errorf("root should not be empty.")
	}
	if !strings.Contains(root, "gira") {
		t.Errorf("Unexpected root: [%s].", root)
	}

	err := os.Chdir(root)
	if err != nil {
		t.Errorf("chdir: %v", err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Getwd: %v", err)
	}
	t.Logf("pwd: %s\n", pwd)
	orig := GetRemote("origin")
	if orig == "" {
		t.Errorf("Unexpected remote: %s.", orig)
	}
}
