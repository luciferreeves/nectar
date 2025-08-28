package screens

import (
	"nectar/components/root"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type rootScreen struct {
	mainArea root.MainAreaModel
}

func _root() tea.Model {
	return &rootScreen{
		mainArea: root.NewMainArea(),
	}
}

func (r *rootScreen) Init() tea.Cmd {
	return r.mainArea.Init()
}

func (r *rootScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		globals.Width, globals.Height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return r, tea.Quit
		}
	}

	var cmd tea.Cmd
	r.mainArea, cmd = r.mainArea.Update(msg)
	return r, cmd
}

func (r *rootScreen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			root.Sidebar(&globals),
			root.MainArea(&globals, r.mainArea),
		),
		root.StatusBar(&globals),
	)
}
