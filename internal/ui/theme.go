package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name string

	Primary   lipgloss.Color
	PrimaryBg lipgloss.Color

	Muted lipgloss.Color
	Text  lipgloss.Color

	Success   lipgloss.Color
	SuccessBg lipgloss.Color
	Warning   lipgloss.Color
	WarningBg lipgloss.Color
	Error     lipgloss.Color
	ErrorBg   lipgloss.Color

	BorderPrimary lipgloss.Color
	BorderSuccess lipgloss.Color
	BorderWarning lipgloss.Color
	BorderError   lipgloss.Color
}

var themes = map[string]*Theme{
	"catppuccin": {
		Name:          "catppuccin",
		Primary:       lipgloss.Color("#89B4FA"),
		PrimaryBg:     lipgloss.Color("#11111B"),
		Muted:         lipgloss.Color("#6C7086"),
		Text:          lipgloss.Color("#CDD6F4"),
		Success:       lipgloss.Color("#A6E3A1"),
		SuccessBg:     lipgloss.Color("#11111B"),
		Warning:       lipgloss.Color("#F9E2AF"),
		WarningBg:     lipgloss.Color("#11111B"),
		Error:         lipgloss.Color("#F38BA8"),
		ErrorBg:       lipgloss.Color("#11111B"),
		BorderPrimary: lipgloss.Color("#585B70"),
		BorderSuccess: lipgloss.Color("#A6E3A1"),
		BorderWarning: lipgloss.Color("#F9E2AF"),
		BorderError:   lipgloss.Color("#F38BA8"),
	},
	"minimal": {
		Name:          "minimal",
		Primary:       lipgloss.Color("7"),
		PrimaryBg:     lipgloss.Color("0"),
		Muted:         lipgloss.Color("8"),
		Text:          lipgloss.Color("7"),
		Success:       lipgloss.Color("2"),
		SuccessBg:     lipgloss.Color("0"),
		Warning:       lipgloss.Color("3"),
		WarningBg:     lipgloss.Color("0"),
		Error:         lipgloss.Color("1"),
		ErrorBg:       lipgloss.Color("0"),
		BorderPrimary: lipgloss.Color("8"),
		BorderSuccess: lipgloss.Color("2"),
		BorderWarning: lipgloss.Color("3"),
		BorderError:   lipgloss.Color("1"),
	},
}

var currentThemeName = "catppuccin"

func SetTheme(name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return false
	}
	if _, ok := themes[name]; !ok {
		return false
	}
	currentThemeName = name
	return true
}

func GetCurrentTheme() *Theme {
	if t, ok := themes[currentThemeName]; ok {
		return t
	}
	return themes["catppuccin"]
}

func GetTheme(name string) (*Theme, bool) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, false
	}
	t, ok := themes[name]
	return t, ok
}

func GetAvailableThemes() []string {
	return []string{currentThemeName}
}
