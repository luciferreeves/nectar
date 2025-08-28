package shared

import (
	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

type ColorOption struct {
	Name  string
	Color lipgloss.AdaptiveColor
}

var ConnectionColors = []ColorOption{
	{
		Name: "Red",
		Color: lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Red().Hex,
			Dark:  catppuccin.Mocha.Red().Hex,
		},
	},
	{
		Name: "Green",
		Color: lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Green().Hex,
			Dark:  catppuccin.Mocha.Green().Hex,
		},
	},
	{
		Name: "Blue",
		Color: lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Blue().Hex,
			Dark:  catppuccin.Mocha.Blue().Hex,
		},
	},
	{
		Name: "Yellow",
		Color: lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Yellow().Hex,
			Dark:  catppuccin.Mocha.Yellow().Hex,
		},
	},
	{
		Name: "Mauve",
		Color: lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Mauve().Hex,
			Dark:  catppuccin.Mocha.Mauve().Hex,
		},
	},
	{
		Name: "Teal",
		Color: lipgloss.AdaptiveColor{
			Light: catppuccin.Latte.Teal().Hex,
			Dark:  catppuccin.Mocha.Teal().Hex,
		},
	},
}
