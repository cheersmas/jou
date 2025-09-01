package navigation

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cheersmas/jou/app/constants"
	"github.com/cheersmas/jou/app/models"
)

type Router struct {
	state *AppState
}

func NewRouter(state *AppState) *Router {
	return &Router{state: state}
}

func (r *Router) LoadJournals() error {
	journals, err := r.state.Service.ListAll(r.state.Ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch journals: %w", err)
	}

	r.state.Journals = journals
	var items []list.Item
	for _, journal := range journals {
		items = append(items, models.NewJournalItem(journal))
	}

	r.state.List.SetItems(items)
	r.state.List.Title = "Journals"
	return nil
}

func (r *Router) HandleMenuSelection() tea.Cmd {
	selectedView := r.state.Options[r.state.CursorPosition]
	r.state.CurrentView = selectedView
	r.state.ResetCursorPosition()

	if selectedView == constants.AddView {
		r.state.Textarea.Focus()
	}

	if selectedView == constants.ListView || selectedView == constants.EditView {
		if err := r.LoadJournals(); err != nil {
			r.state.LastError = err
			log.Printf("Error loading journals: %v", err)
		}
	}
	return nil
}

func (r *Router) HandleJournalSelection() tea.Cmd {
	if r.state.CursorPosition >= len(r.state.Journals) {
		return nil
	}

	selected := r.state.Journals[r.state.CursorPosition]

	if r.state.CurrentView == constants.EditView {
		r.state.EditingJournal = &selected
		r.state.CurrentView = constants.AddView
		r.state.RecentlySavedId = selected.Id
		content := strings.TrimSpace(selected.Content)
		r.state.Textarea.SetValue(content)
	} else {
		r.state.ViewingJournal = &selected
		r.state.Viewport.SetContent(selected.Content)
		r.state.CurrentView = constants.JournalView
	}

	r.state.ResetCursorPosition()
	return nil
}

func (r *Router) HasUnsavedChanges() bool {
	if r.state.RecentlySavedId == constants.UnsavedId {
		return true
	}
	return r.state.RecentlySavedId != constants.UnsavedId &&
		r.state.Textarea.Value() != r.state.EditingJournal.Content
}
