package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cheersmas/jou/app/navigation"
)

type View interface {
	Render(state *navigation.AppState) string
	Update(state *navigation.AppState, msg tea.Msg) tea.Cmd
}
