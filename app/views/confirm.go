package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/app/styles"
)

type ConfirmView struct{}

func (v ConfirmView) Render(state *navigation.AppState) string {
	header := styles.HeaderStyle.Render("Confirm Exit")

	content := "Unsaved changes may get lost\n\n"
	content += "• <esc>, <backspace>: cancel\n"
	content += "• <enter>: discard and go to main menu\n"
	content += "• <ctrl + q>: quit"

	footer := styles.FooterStyle.Render("Choose an option above")

	return styles.ContainerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, header, "", content, "", footer),
	)
}

func (v ConfirmView) Update(state *navigation.AppState, msg tea.Msg) tea.Cmd {
	return nil
}
