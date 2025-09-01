package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cheersmas/jou/app/constants"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/app/styles"
)

type JournalView struct{}

func (v JournalView) Render(state *navigation.AppState) string {
	if !state.Ready {
		return "\n  Initializing..."
	}

	header := v.headerView(state)
	content := state.Viewport.View()
	footer := v.footerView(state)

	return fmt.Sprintf("%s\n%s\n%s", header, content, footer)
}

func (v JournalView) Update(state *navigation.AppState, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	state.Viewport, cmd = state.Viewport.Update(msg)
	return cmd
}

func (v JournalView) headerView(state *navigation.AppState) string {
	createdAt := "Untitled"
	if state.ViewingJournal != nil {
		createdAt = state.ViewingJournal.CreatedAt.Format(constants.TimeFormat)
	}
	title := styles.TitleStyle.Render(createdAt)
	line := strings.Repeat("─", max(0, state.Viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (v JournalView) footerView(state *navigation.AppState) string {
	// Create the navigation footer on the left
	navFooter := styles.FooterStyle.Render("↑k up • ↓j down • esc back to list • ctrl+c quit")

	// Create the scroll percentage on the right
	scrollInfo := styles.InfoStyle.Render(fmt.Sprintf("%3.f%%", state.Viewport.ScrollPercent()*100))

	// Calculate the space needed, ensuring it's not negative
	navWidth := lipgloss.Width(navFooter)
	scrollWidth := lipgloss.Width(scrollInfo)
	availableWidth := state.Viewport.Width - navWidth - scrollWidth

	// If there's not enough space, just put them next to each other
	if availableWidth <= 0 {
		return lipgloss.JoinHorizontal(lipgloss.Left, navFooter, scrollInfo)
	}

	// Join them side by side with space between
	return lipgloss.JoinHorizontal(lipgloss.Left, navFooter, strings.Repeat(" ", availableWidth), scrollInfo)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
