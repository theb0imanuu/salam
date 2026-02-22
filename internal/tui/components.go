package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	subtle    = lipgloss.AdaptiveColor{Light: "#D9D9D9", Dark: "#383838"}
	success   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warning   = lipgloss.AdaptiveColor{Light: "#FFA500", Dark: "#FFD700"}
	danger    = lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#FF5F87"}

	// Styles
	mainStyle = lipgloss.NewStyle().
			Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(1).
			MarginRight(1).
			MarginBottom(1)

	metricLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	metricValueStyle = lipgloss.NewStyle().
				Bold(true)
)

func renderProgressBar(width int, percentage float64) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}

	filledLen := int(float64(width) * (percentage / 100))
	emptyLen := width - filledLen

	barColor := success
	if percentage > 85 {
		barColor = danger
	} else if percentage > 70 {
		barColor = warning
	}

	filled := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", filledLen))
	empty := lipgloss.NewStyle().Foreground(subtle).Render(strings.Repeat("░", emptyLen))

	return fmt.Sprintf("[%s%s] %.1f%%", filled, empty, percentage)
}

func renderBox(title string, content string) string {
	return boxStyle.Render(
		fmt.Sprintf("%s\n\n%s",
			lipgloss.NewStyle().Bold(true).Foreground(highlight).Render(title),
			content,
		),
	)
}
