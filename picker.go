package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.
//
// Thanks to Treilik for this one.

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wzhliang/gira/pkg/context"
)

type Lister interface {
	List(ctx *context.Context) ([]string, error)
}

var choices []string

type model struct {
	cursor int
	choice chan string
	prompt string
}

func Pick(cmd *CmdContext, l Lister, prompt string) string {
	result := make(chan string, 1)

	var err error
	choices, err = l.List(cmd.ctx)
	if err != nil {
		return ""
	}

	p := tea.NewProgram(model{
		cursor: 0,
		choice: result,
		prompt: prompt,
	})
	if err := p.Start(); err != nil {
		fmt.Println("Oh no:", err)
		return ""
	}

	if r := <-result; r != "" {
		// TODO: this assumes the format, should be implemented by
		// individual providers like jira, or gitee
		// Also, Picker should not knwo this
		return strings.Split(r, " ")[0]
	}

	return ""
}

func initialModel(choice chan string) model {
	return model{cursor: 0, choice: choice}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			close(m.choice)
			return m, tea.Quit

		case "enter":
			m.choice <- choices[m.cursor]
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString(m.prompt)
	s.WriteString("\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("[*] ")
		} else {
			s.WriteString("[ ] ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n")

	return s.String()
}
