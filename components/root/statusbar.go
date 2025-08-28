package root

import (
	"nectar/styles"
	"nectar/types"

	"github.com/charmbracelet/lipgloss"
)

func StatusBar(globals *types.Globals) string {
	w := lipgloss.Width

	helpText := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styles.PaddedHorizontal.Render("↑/k: up"),
		styles.PaddedHorizontal.Render("↓/j: down"),
		styles.PaddedHorizontal.Render("↹: next"),
		styles.PaddedHorizontal.Render("^n: new"),
		styles.PaddedHorizontal.Render("^↵: connect"),
		styles.PaddedHorizontal.Render("^s: save"),
		styles.PaddedHorizontal.Render("^t: test"),
		styles.PaddedHorizontal.Render("^d: delete"),
		styles.PaddedHorizontal.Render("^c: quit"),
	)

	versionText := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styles.PaddedHorizontal.Render("Nectar "+globals.Version+" ("+globals.BuildDate+")"),
	)

	separator := styles.BaseStyle.Width(globals.Width - w(helpText) - w(versionText)).Render("")

	return styles.StatusBar.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			helpText,
			separator,
			versionText,
		),
	)
}
