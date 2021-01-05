package main

import (
	"testing"

	"github.com/wzhliang/gira/pkg/context"
)

// Mergable ...
type Mergable struct {
	msg string
}

// Check ...
func (jc *Mergable) Check(ctx *context.Context) bool {
	jc.msg = "Unknown error."
	return false
}

// Message ...
func (jc *Mergable) Message(ctx *context.Context) string {
	return "Checking if JIRA issue is mergable..."
}

// ErrMessage ...
func (jc *Mergable) ErrMessage(ctx *context.Context) string {
	return jc.msg
}

func TestPolicy(t *testing.T) {

	m := Mergable{}
	m.Check(nil)
	t.Logf("--- %s\n", m.Message(nil))
	t.Logf("--- %s\n", m.ErrMessage(nil))

	pol := &Policy{}
	if pol.
		Add(Enforcer(&Mergable{})).
		Check(&_ctx) == true {
		t.Errorf("should have failed.")
	}
}
