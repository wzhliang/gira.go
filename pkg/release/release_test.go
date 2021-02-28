package release

import "testing"

func TestBasicRelease(t *testing.T) {
	rls := New("v1.12.2")
	if rls == nil {
		t.Errorf("Should have parsed OK")
		return
	}
	if rls.Major != 1 {
		t.Errorf("expect major == 1")
	}

	rls = New("v2.0.0-foo")
	if rls == nil {
		t.Errorf("Should have parsed OK")
		return
	}
	if rls.Major != 2 {
		t.Errorf("expect major == 0")
	}
}

func TestReleaseBranch(t *testing.T) {
	rls := New("v1.12.2")

	if rls.TargetBranch() != "release-1.12" {
		t.Errorf("expect target branch == release-1.12")
	}

	rls = New("v1.12.0")
	if rls.TargetBranch() != "master" {
		t.Errorf("expect target branch == master")
	}

	rls = New("v2.0.0")
	if rls.TargetBranch() != "master" {
		t.Errorf("expect target branch == master")
	}
}
