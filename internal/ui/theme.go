package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name string

	Primary   lipgloss.Color
	PrimaryBg lipgloss.Color
	Secondary lipgloss.Color

	Muted lipgloss.Color
	Text  lipgloss.Color

	Success   lipgloss.Color
	SuccessBg lipgloss.Color
	Warning   lipgloss.Color
	WarningBg lipgloss.Color
	Error     lipgloss.Color
	ErrorBg   lipgloss.Color

	BorderPrimary   lipgloss.Color
	BorderSecondary lipgloss.Color
	BorderSuccess   lipgloss.Color
	BorderWarning   lipgloss.Color
	BorderError     lipgloss.Color

	Accent lipgloss.Color
	Dim    lipgloss.Color
}

var themes = map[string]*Theme{
	"catppuccin": {
		Name:            "catppuccin",
		Primary:         lipgloss.Color("#89B4FA"),
		PrimaryBg:       lipgloss.Color("#11111B"),
		Secondary:       lipgloss.Color("#CBA6F7"),
		Muted:           lipgloss.Color("#6C7086"),
		Text:            lipgloss.Color("#CDD6F4"),
		Success:         lipgloss.Color("#A6E3A1"),
		SuccessBg:       lipgloss.Color("#11111B"),
		Warning:         lipgloss.Color("#F9E2AF"),
		WarningBg:       lipgloss.Color("#11111B"),
		Error:           lipgloss.Color("#F38BA8"),
		ErrorBg:         lipgloss.Color("#11111B"),
		BorderPrimary:   lipgloss.Color("#45475A"),
		BorderSecondary: lipgloss.Color("#9399B2"),
		BorderSuccess:   lipgloss.Color("#A6E3A1"),
		BorderWarning:   lipgloss.Color("#F9E2AF"),
		BorderError:     lipgloss.Color("#F38BA8"),
		Accent:          lipgloss.Color("#F5C2E7"),
		Dim:             lipgloss.Color("#313244"),
	},
	"night": {
		Name:            "night",
		Primary:         lipgloss.Color("#7AA2F7"),
		PrimaryBg:       lipgloss.Color("#1A1B26"),
		Secondary:       lipgloss.Color("#BB9AF7"),
		Muted:           lipgloss.Color("#565F89"),
		Text:            lipgloss.Color("#C0CAF5"),
		Success:         lipgloss.Color("#9ECE6A"),
		SuccessBg:       lipgloss.Color("#1A1B26"),
		Warning:         lipgloss.Color("#E0AF68"),
		WarningBg:       lipgloss.Color("#1A1B26"),
		Error:           lipgloss.Color("#F7768E"),
		ErrorBg:         lipgloss.Color("#1A1B26"),
		BorderPrimary:   lipgloss.Color("#3B4261"),
		BorderSecondary: lipgloss.Color("#7AA2F7"),
		BorderSuccess:   lipgloss.Color("#9ECE6A"),
		BorderWarning:   lipgloss.Color("#E0AF68"),
		BorderError:     lipgloss.Color("#F7768E"),
		Accent:          lipgloss.Color("#FF9E64"),
		Dim:             lipgloss.Color("#292E42"),
	},
	"minimal": {
		Name:            "minimal",
		Primary:         lipgloss.Color("7"),
		PrimaryBg:       lipgloss.Color("0"),
		Secondary:       lipgloss.Color("7"),
		Muted:           lipgloss.Color("8"),
		Text:            lipgloss.Color("7"),
		Success:         lipgloss.Color("2"),
		SuccessBg:       lipgloss.Color("0"),
		Warning:         lipgloss.Color("3"),
		WarningBg:       lipgloss.Color("0"),
		Error:           lipgloss.Color("1"),
		ErrorBg:         lipgloss.Color("0"),
		BorderPrimary:   lipgloss.Color("8"),
		BorderSecondary: lipgloss.Color("7"),
		BorderSuccess:   lipgloss.Color("2"),
		BorderWarning:   lipgloss.Color("3"),
		BorderError:     lipgloss.Color("1"),
		Accent:          lipgloss.Color("6"),
		Dim:             lipgloss.Color("0"),
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

func GetAvailableThemes() []string {
	names := make([]string, 0, len(themes))
	for k := range themes {
		names = append(names, k)
	}
	return names
}
