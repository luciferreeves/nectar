package screens

import (
	"nectar/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type rootScreen struct{}

func _root() tea.Model {
	return &rootScreen{}
}

func (r *rootScreen) Init() tea.Cmd {
	return nil
}

func (r *rootScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return r, tea.Quit
		case "a":
			return r, switchScreen(_aux())
		}
	}
	return r, nil
}

func (r *rootScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		components.RootStatusBar(&globals),
	)
}
