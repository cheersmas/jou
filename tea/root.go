package tea

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cheersmas/jou/domains"
	"github.com/cheersmas/jou/ports"
)

type View string

const (
	MenuView    View = "Menu"
	AddView     View = "Add"
	ListView    View = "View"
	EditView    View = "Edit"
	ConfirmView View = "Confirm"

	timeFormat = "2 Jan, 2006"
)

// Constants for state management
const (
	unsavedId = -1
)

type model struct {
	ctx     context.Context
	service ports.JournalService

	// Navigation state
	currentView    View
	options        []View
	cursorPosition int

	// Journal data
	journals       []domains.Journal
	editingJournal *domains.Journal // Use pointer for clear nil state

	// Text area state
	textarea        textarea.Model
	recentlySavedId int
	lastError       error
}

func initialModel(ctx context.Context, js ports.JournalService) model {
	ti := textarea.New()
	ti.Placeholder = "Write your journal entry here..."
	ti.Focus()

	return model{
		ctx:             ctx,
		service:         js,
		options:         []View{AddView, ListView, EditView},
		currentView:     MenuView,
		textarea:        ti,
		recentlySavedId: unsavedId,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// loadJournals loads journals and returns updated model
func (m *model) loadJournals() error {
	journals, err := m.service.ListAll(m.ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch journals: %w", err)
	}
	m.journals = journals
	return nil
}

// handleKeyMsg processes all keyboard input
func (m *model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		return m.handleEscapeKey()
	case "ctrl+c", "q":
		return m.handleQuitKey()
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		return m.handleEnterKey()
	case "ctrl+s":
		return m.handleSaveKey()
	case "backspace":
		return m.handleBackspaceKey()
	default:
		return m.handleDefaultKey()
	}
	return nil
}

// Navigation handlers
func (m *model) handleEscapeKey() tea.Cmd {
	switch m.currentView {
	case AddView:
		m.textarea.Blur()
		return nil
	case ConfirmView:
		m.currentView = AddView
		return nil
	}
	return nil
}

func (m *model) handleQuitKey() tea.Cmd {
	switch m.currentView {
	case AddView:
		m.currentView = ConfirmView
		return nil
	default:
		return tea.Quit
	}
}

func (m *model) moveCursor(direction int) {
	switch m.currentView {
	case MenuView:
		newPos := m.cursorPosition + direction
		if newPos >= 0 && newPos < len(m.options) {
			m.cursorPosition = newPos
		}
	case ListView, EditView:
		newPos := m.cursorPosition + direction
		if newPos >= 0 && newPos < len(m.journals) {
			m.cursorPosition = newPos
		}
	}
}

func (m *model) handleEnterKey() tea.Cmd {
	switch m.currentView {
	case MenuView:
		return m.handleMenuSelection()
	case ListView, EditView:
		return m.handleJournalSelection()
	case ConfirmView:
		m.resetCursorPosition()
		m.currentView = MenuView
		return nil
	}
	return nil
}

func (m *model) resetCursorPosition() {
	m.cursorPosition = 0
}

func (m *model) handleMenuSelection() tea.Cmd {
	selectedView := m.options[m.cursorPosition]
	m.currentView = selectedView
	m.resetCursorPosition()

	if selectedView == ListView || selectedView == EditView {
		if err := m.loadJournals(); err != nil {
			m.lastError = err
			log.Printf("Error loading journals: %v", err)
		}
	}
	return nil
}

func (m *model) handleJournalSelection() tea.Cmd {
	if m.cursorPosition >= len(m.journals) {
		return nil
	}

	selected := m.journals[m.cursorPosition]

	if m.currentView == EditView {
		m.editingJournal = &selected
		m.textarea.SetValue(selected.Content)
		m.currentView = AddView
		m.recentlySavedId = selected.Id
	} else {
		// Handle view mode - could be extracted to separate view
		m.editingJournal = &selected
	}

	m.resetCursorPosition()
	return nil
}

func (m *model) handleSaveKey() tea.Cmd {
	if m.currentView != AddView {
		return nil
	}

	content := m.textarea.Value()
	if content == "" {
		return nil
	}

	var err error
	if m.recentlySavedId == unsavedId {
		journal := domains.Journal{Content: content}
		m.recentlySavedId, err = m.service.Create(m.ctx, journal)
	} else {
		val, err := m.service.Update(m.ctx, m.recentlySavedId, content)
		fmt.Printf("value %d", val)
		if err != nil {
			m.lastError = err
			log.Printf("Save error: %v", err)
		}
		m.recentlySavedId = val
	}

	if err != nil {
		m.lastError = err
		log.Printf("Save error: %v", err)
	}
	return nil
}

func (m *model) handleBackspaceKey() tea.Cmd {
	if m.currentView == ListView || m.currentView == EditView {
		m.currentView = MenuView
		m.cursorPosition = 0
	}
	return nil
}

func (m *model) handleDefaultKey() tea.Cmd {
	if m.currentView == AddView && !m.textarea.Focused() {
		return m.textarea.Focus()
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = m.handleKeyMsg(msg)
	case error:
		m.lastError = msg
		return m, nil
	}

	// Handle textarea updates only in Add view
	if m.currentView == AddView {
		m.textarea, cmd = m.textarea.Update(msg)
	}

	return m, cmd
}

// View rendering methods
func (m model) renderMenu() string {
	s := "What would you like to do?\n\n"
	for i, option := range m.options {
		cursor := "[ ]"
		if i == m.cursorPosition {
			cursor = "[>]"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}
	return s
}

func (m model) renderConfirmation() string {
	return "Unsaved changes will be lost.\n\nPress:\n<esc> to cancel\n<enter> to discard"
}

func (m model) renderAddJournal() string {
	var createdAt time.Time
	if m.editingJournal != nil {
		createdAt = m.editingJournal.CreatedAt
	} else {
		createdAt = time.Now()
	}

	s := fmt.Sprintf("Journal entry %s\n\n", createdAt.Format(timeFormat))
	s += m.textarea.View()

	if m.recentlySavedId != unsavedId {
		s += fmt.Sprintf("\nSaved (ID: %d)\n", m.recentlySavedId)
	}

	if m.lastError != nil {
		s += fmt.Sprintf("\nError: %v\n", m.lastError)
	}

	return s
}

func (m model) renderJournalList() string {
	if len(m.journals) == 0 {
		return "No journals found.\n\nPress backspace to return to menu"
	}

	s := "All Journals\n\n"
	for i, journal := range m.journals {
		cursor := "[ ]"
		if i == m.cursorPosition {
			cursor = "[>]"
		}

		// Truncate content for display
		contentPreview := journal.Content
		if len(contentPreview) > 50 {
			contentPreview = contentPreview[:47] + "..."
		}

		formattedDate := journal.CreatedAt.Format(timeFormat)
		s += fmt.Sprintf("%s %s | %s\n", cursor, formattedDate, contentPreview)
	}
	return s
}

func (m model) View() string {
	var s string

	switch m.currentView {
	case MenuView:
		s = m.renderMenu()
	case AddView:
		s = m.renderAddJournal()
	case ListView, EditView:
		s = m.renderJournalList()
	case ConfirmView:
		s = m.renderConfirmation()
	}

	if m.currentView != ConfirmView {
		s += "\n(Press ctrl+c or q to quit)\n"
	}

	return s
}

func Root(ctx context.Context, js ports.JournalService) {
	p := tea.NewProgram(initialModel(ctx, js))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
