package styles

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseStyle        = lipgloss.NewStyle()
	PaddedHorizontal = BaseStyle.Padding(0, 1)
	StatusBar        = BaseStyle.
				Foreground(lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Crust().Hex,
			Dark:  catppuccin.Mocha.Crust().Hex,
		}).
		Background(lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Mauve().Hex,
			Dark:  catppuccin.Mocha.Mauve().Hex,
		})
)
