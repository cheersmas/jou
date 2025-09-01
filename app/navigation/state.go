package navigation

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/cheersmas/jou/app/constants"
	"github.com/cheersmas/jou/domains"
	"github.com/cheersmas/jou/ports"
)

type AppState struct {
	// Core dependencies
	Ctx     context.Context
	Service ports.JournalService

	// Navigation state
	CurrentView    constants.View
	Options        []constants.View
	CursorPosition int

	// Journal data
	Journals       []domains.Journal
	List           list.Model
	ViewingJournal *domains.Journal
	EditingJournal *domains.Journal

	// UI components
	Viewport        viewport.Model
	Textarea        textarea.Model
	RecentlySavedId int
	LastError       error
	Ready           bool
}

func NewAppState(ctx context.Context, service ports.JournalService) *AppState {
	return &AppState{
		Ctx:             ctx,
		Service:         service,
		Options:         []constants.View{constants.AddView, constants.ListView, constants.EditView},
		CurrentView:     constants.MenuView,
		RecentlySavedId: constants.UnsavedId,
		Ready:           false,
	}
}

func (s *AppState) ResetCursorPosition() {
	s.CursorPosition = 0
}

func (s *AppState) MoveCursor(direction int) {
	switch s.CurrentView {
	case constants.MenuView:
		newPos := s.CursorPosition + direction
		if newPos >= 0 && newPos < len(s.Options) {
			s.CursorPosition = newPos
		}
		// Remove the ListView and EditView cases since the list component handles its own cursor
	}
}
