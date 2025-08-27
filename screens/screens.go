package screens

import (
	"nectar/build"
	"nectar/types"
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

var globals types.Globals

func (sm *ScreenManager) Init() tea.Cmd {
	return sm.currentScreen.Init()
}

func (sm *ScreenManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		globals.Width, globals.Height = m.Width, m.Height
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

	globals = types.Globals{
		Width:     width,
		Height:    height,
		BuildDate: build.Date,
		Version:   build.Version,
	}

	return &ScreenManager{currentScreen: _root()}
}
