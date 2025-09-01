package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/app/styles"
)

type MenuView struct{}

func (v MenuView) Render(state *navigation.AppState) string {
	header := styles.HeaderStyle.Render("jou")
	subtitle := styles.FooterStyle.Render("A commandline journaling tool")

	content := "What would you like to do?\n\n"
	for i, option := range state.Options {
		cursor := "[ ]"
		if i == state.CursorPosition {
			cursor = "[>]"
		}
		content += fmt.Sprintf("%s %s\n", cursor, option)
	}

	footer := styles.FooterStyle.Render("↑k up • ↓j down • enter select • ctrl+c quit")

	fullContent := lipgloss.JoinVertical(lipgloss.Left, header, subtitle, "", content, "", footer)

	return lipgloss.NewStyle().
		Width(state.Viewport.Width).
		Height(state.Viewport.Height).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(fullContent)
}

func (v MenuView) Update(state *navigation.AppState, msg tea.Msg) tea.Cmd {
	return nil
}
