package input

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cheersmas/jou/app/constants"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/domains"
)

type InputHandler struct {
	state  *navigation.AppState
	router *navigation.Router
}

func NewInputHandler(state *navigation.AppState, router *navigation.Router) *InputHandler {
	return &InputHandler{
		state:  state,
		router: router,
	}
}

func (h *InputHandler) HandleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		return h.handleEscapeKey()
	case "ctrl+c":
		return h.handleQuitKey(msg)
	case "up", "k":
		// Only move cursor for menu view, let list component handle its own navigation
		if h.state.CurrentView == constants.MenuView {
			h.state.MoveCursor(-1)
		}
	case "down", "j":
		// Only move cursor for menu view, let list component handle its own navigation
		if h.state.CurrentView == constants.MenuView {
			h.state.MoveCursor(1)
		}
	case "enter":
		return h.handleEnterKey()
	case "ctrl+s":
		return h.handleSaveKey()
	case "backspace":
		return h.handleBackspaceKey()
	default:
		return h.handleDefaultKey()
	}
	return nil
}

func (h *InputHandler) handleEscapeKey() tea.Cmd {
	switch h.state.CurrentView {
	case constants.AddView:
		h.state.Textarea.Blur()
		return nil
	case constants.ConfirmView:
		h.state.CurrentView = constants.AddView
		return nil
	}
	return nil
}

func (h *InputHandler) handleQuitKey(msg tea.KeyMsg) tea.Cmd {
	switch h.state.CurrentView {
	case constants.AddView:
		h.state.CurrentView = constants.ConfirmView
		return nil
	default:
		return tea.Quit
	}
}

func (h *InputHandler) handleEnterKey() tea.Cmd {
	switch h.state.CurrentView {
	case constants.MenuView:
		return h.router.HandleMenuSelection()
	case constants.ListView, constants.EditView:
		return h.router.HandleJournalSelection()
	case constants.ConfirmView:
		h.state.ResetCursorPosition()
		h.state.CurrentView = constants.MenuView
		return nil
	}
	return nil
}

func (h *InputHandler) handleSaveKey() tea.Cmd {
	if h.state.CurrentView != constants.AddView {
		return nil
	}

	content := h.state.Textarea.Value()
	if content == "" {
		return nil
	}

	var err error
	if h.state.RecentlySavedId == constants.UnsavedId {
		journal := domains.Journal{Content: content}
		h.state.RecentlySavedId, err = h.state.Service.Create(h.state.Ctx, journal)
	} else {
		val, err := h.state.Service.Update(h.state.Ctx, h.state.RecentlySavedId, content)
		if err != nil {
			h.state.LastError = err
			log.Printf("Save error: %v", err)
		}
		h.state.RecentlySavedId = val
	}

	editingJournal, err := h.state.Service.Read(h.state.Ctx, h.state.RecentlySavedId)
	h.state.EditingJournal = &editingJournal

	if err != nil {
		h.state.LastError = err
		log.Printf("Save error: %v", err)
	}
	return nil
}

func (h *InputHandler) handleBackspaceKey() tea.Cmd {
	switch h.state.CurrentView {
	case constants.ListView, constants.EditView:
		h.state.CurrentView = constants.MenuView
		h.state.CursorPosition = 0
	case constants.JournalView:
		h.state.CurrentView = constants.ListView
	case constants.ConfirmView:
		h.state.CurrentView = constants.AddView
	}
	return nil
}

func (h *InputHandler) handleDefaultKey() tea.Cmd {
	if h.state.CurrentView == constants.AddView && !h.state.Textarea.Focused() {
		return h.state.Textarea.Focus()
	}
	return nil
}
