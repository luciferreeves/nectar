package root

import (
	"nectar/styles"
	"nectar/types"

	"github.com/charmbracelet/lipgloss"
)

func Sidebar(globals *types.Globals) string {
	return styles.BaseStyle.
		Width(30).Height(globals.Height - 1).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		Render("Sidebar placeholder")
}
