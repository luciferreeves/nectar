package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func _aux() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &auxScreen{
		spinner: s,
	}
}

type auxScreen struct {
	spinner spinner.Model
}

func (a *auxScreen) Init() tea.Cmd {
	return a.spinner.Tick
}

func (a *auxScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "r":
			return a, switchScreen(_root())
		default:
			return a, nil
		}
	default:
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		return a, cmd
	}
}

func (a *auxScreen) View() string {
	str := fmt.Sprintf("%s Loading forever... Press 'r' to return to root or 'q' to quit.\n\n", a.spinner.View())
	return str
}
