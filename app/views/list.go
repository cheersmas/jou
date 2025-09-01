package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/app/styles"
)

type ListView struct{}

func (v ListView) Render(state *navigation.AppState) string {
	return styles.DocStyle.Render(state.List.View())
}

func (v ListView) Update(state *navigation.AppState, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	state.List, cmd = state.List.Update(msg)
	return cmd
}
