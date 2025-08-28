package root

import (
	"nectar/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainAreaModel struct {
	connectionForm ConnectionFormModel
}

func NewMainArea() MainAreaModel {
	return MainAreaModel{
		connectionForm: NewConnectionForm(),
	}
}

func (m MainAreaModel) Init() tea.Cmd {
	return m.connectionForm.Init()
}

func (m MainAreaModel) Update(msg tea.Msg) (MainAreaModel, tea.Cmd) {
	var cmd tea.Cmd
	formModel, cmd := m.connectionForm.Update(msg)
	m.connectionForm = formModel.(ConnectionFormModel)
	return m, cmd
}

func MainArea(globals *types.Globals, mainArea MainAreaModel) string {
	formContent := mainArea.connectionForm.View()

	return lipgloss.Place(
		globals.Width-30,
		globals.Height-1,
		lipgloss.Center,
		lipgloss.Center,
		formContent,
	)
}
