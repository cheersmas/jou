package views

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cheersmas/jou/app/constants"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/app/styles"
)

type AddView struct{}

func (v AddView) Render(state *navigation.AppState) string {
	header := v.addJournalHeader(state)
	content := state.Textarea.View()
	footer := v.addJournalFooter(state)

	return styles.ContainerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, header, "", content, "", footer),
	)
}

func (v AddView) Update(state *navigation.AppState, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	state.Textarea, cmd = state.Textarea.Update(msg)
	return cmd
}

func (v AddView) addJournalHeader(state *navigation.AppState) string {
	var createdAt time.Time
	if state.EditingJournal != nil {
		createdAt = state.EditingJournal.CreatedAt
	} else {
		createdAt = time.Now()
	}

	header := styles.HeaderStyle.Render(fmt.Sprintf("Journal entry %s", createdAt.Format("2 Jan, 2006")))
	return header
}

func (v AddView) addJournalFooter(state *navigation.AppState) string {
	var status string
	router := navigation.NewRouter(state)
	if router.HasUnsavedChanges() {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("● Unsaved changes")
	} else {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("✓ Saved")
	}

	if state.RecentlySavedId != constants.UnsavedId {
		status += fmt.Sprintf(" (ID: %d)", state.RecentlySavedId)
	}

	if state.LastError != nil {
		status += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("✗ Error: %v", state.LastError))
	}

	footer := styles.FooterStyle.Render("<ctrl + s>: Save | <ctrl + c>: quit")

	return lipgloss.JoinVertical(lipgloss.Left, status, "", footer)
}
