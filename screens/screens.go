package screens

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

type SwitchMsg struct {
	Screen tea.Model
}

type ScreenManager struct {
	currentScreen tea.Model
}

var globals Globals

func (sm *ScreenManager) Init() tea.Cmd {
	return sm.currentScreen.Init()
}

func (sm *ScreenManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		globals.width, globals.height = m.Width, m.Height
		return sm, nil
	case SwitchMsg:
		sm.currentScreen = m.Screen
		return sm, sm.currentScreen.Init()
	}

	updatedScreen, cmd := sm.currentScreen.Update(msg)
	sm.currentScreen = updatedScreen
	return sm, cmd
}

func (sm *ScreenManager) View() string {
	return sm.currentScreen.View()
}

func switchScreen(screen tea.Model) tea.Cmd {
	return func() tea.Msg {
		return SwitchMsg{Screen: screen}
	}
}

func Start() tea.Model {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
		height = 24
	}

	globals = Globals{width: width, height: height}

	return &ScreenManager{currentScreen: _root()}
}
