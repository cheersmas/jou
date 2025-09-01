package tea

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cheersmas/jou/domains"
	"github.com/cheersmas/jou/ports"
)

type View string

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	roundedBorderStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

const (
	MenuView    View = "Menu"
	AddView     View = "Add"
	ListView    View = "View"
	JournalView View = "Journal"
	EditView    View = "Edit"
	ConfirmView View = "Confirm"

	timeFormat = "2 Jan, 2006"
	gap        = "\n\n"
)

// Constants for state management
const (
	unsavedId = -1
)

type viewportDimensions struct {
	height int
	width  int
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.desc }

type model struct {
	// dimensions
	viewport           viewport.Model
	viewportDimensions viewportDimensions

	ctx     context.Context
	service ports.JournalService

	// Navigation state
	currentView    View
	options        []View
	cursorPosition int

	// Journal data
	journals       []domains.Journal
	list           list.Model
	viewingJournal *domains.Journal
	editingJournal *domains.Journal // Use pointer for clear nil state

	// Text area state
	textarea        textarea.Model
	recentlySavedId int
	lastError       error

	ready bool
}

// Add the missing max function
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func initialModel(ctx context.Context, js ports.JournalService) model {
	ti := textarea.New()
	ti.Placeholder = "Write your journal entry here..."
	ti.Focus()

	vp := viewport.New(30, 5)
	vp.SetContent(`init`)

	items := []list.Item{}
	li := list.New(items, list.NewDefaultDelegate(), 0, 0)

	return model{
		ctx:             ctx,
		service:         js,
		options:         []View{AddView, ListView, EditView},
		currentView:     MenuView,
		textarea:        ti,
		recentlySavedId: unsavedId,
		viewport:        vp,
		ready:           false,
		list:            li,
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
	var items []list.Item
	for _, journal := range journals {
		title := journal.CreatedAt.Format(timeFormat)
		items = append(items, item{title: title, desc: journal.Content})
	}
	m.list = list.New(items, list.NewDefaultDelegate(), m.list.Width(), m.list.Height())
	m.list.Title = "Journals"
	return nil
}

// handleKeyMsg processes all keyboard input
func (m *model) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		return m.handleEscapeKey()
	case "ctrl+c", "q":
		return m.handleQuitKey(msg)
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

func (m *model) handleQuitKey(msg tea.KeyMsg) tea.Cmd {
	switch m.currentView {
	case AddView:
		if msg.String() != "q" {
			m.currentView = ConfirmView
		}
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
		// set the content first
		m.viewingJournal = &selected
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Height(m.viewport.Height).Render(selected.Content))
		m.currentView = JournalView
	}

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
		if err != nil {
			m.lastError = err
			log.Printf("Save error: %v", err)
		}
		m.recentlySavedId = val
	}

	// TODO: yuck fix this by getting the journal before hand
	m.editingJournal = &domains.Journal{Content: content}

	if err != nil {
		m.lastError = err
		log.Printf("Save error: %v", err)
	}
	return nil
}

func (m *model) handleBackspaceKey() tea.Cmd {
	switch m.currentView {
	case ListView, EditView:
		m.currentView = MenuView
		m.cursorPosition = 0
	case JournalView:
		m.currentView = ListView
	case ConfirmView:
		m.currentView = AddView
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
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = m.handleKeyMsg(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		// list
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		addJournalHeaderHeight := lipgloss.Height(m.addJournalHeader())
		addJournalFooterHeight := lipgloss.Height(m.addJournalFooter())
		addJournalverticalMarginHeight := addJournalFooterHeight + addJournalHeaderHeight
		m.textarea.SetHeight(msg.Height - addJournalverticalMarginHeight - lipgloss.Height(gap))
		m.textarea.SetWidth(msg.Width)

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case error:
		m.lastError = msg
		return m, nil
	}

	switch m.currentView {
	case AddView:
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	case JournalView:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	case ListView, EditView:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View rendering methods
func (m model) renderMenu() string {
	title := roundedBorderStyle.Render("CLJour")
	description := roundedBorderStyle.Render("A commandline journaling tool")
	question := "\nWhat would you like to do?\n\n"
	formatedString := lipgloss.JoinVertical(lipgloss.Center, title, description, question)
	s := ""
	for i, option := range m.options {
		cursor := "[ ]"
		if i == m.cursorPosition {
			cursor = "[>]"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}
	return formatedString + "\n" + s
}

func (m model) renderConfirmation() string {
	confirmationMessage := "Unsaved chages may get lost"
	cancel := roundedBorderStyle.Render("<esc>,<backspace>: cancel")
	discard := roundedBorderStyle.Render("<enter>: discard and go to main menu")
	quit := roundedBorderStyle.Render("<ctrl + q>: quit")
	return lipgloss.JoinVertical(lipgloss.Left, confirmationMessage, cancel, discard, quit)
}

func (m model) hasUnsavedChanges() bool {
	if m.recentlySavedId == unsavedId {
		return true
	}
	return m.recentlySavedId != unsavedId && m.textarea.Value() != m.editingJournal.Content
}

func (m model) addJournalHeader() string {
	var createdAt time.Time
	if m.editingJournal != nil {
		createdAt = m.editingJournal.CreatedAt
	} else {
		createdAt = time.Now()
	}
	s := titleStyle.Render(fmt.Sprintf("Journal entry %s", createdAt.Format(timeFormat)))
	return lipgloss.JoinHorizontal(lipgloss.Center, s)
}

func (m model) addJournalFooter() string {
	info := infoStyle.Render("<ctrl + s>: Save | <ctrl + c>: quit")
	status := ""
	if m.hasUnsavedChanges() {
		status = infoStyle.Render("unsaved changes")
	} else {
		status = infoStyle.Render("saved")
	}
	if m.lastError != nil {
		info = infoStyle.Render(fmt.Sprintf("Error: %v\n", m.lastError))
	}

	line := strings.Repeat("─", max(0, m.viewport.Width-(lipgloss.Width(info)+lipgloss.Width(status))))
	return lipgloss.JoinHorizontal(lipgloss.Center, info, status, line)
}

func (m model) renderAddJournal() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", m.addJournalHeader(), m.textarea.View(), m.addJournalFooter())
}

func (m model) renderJournalList() string {
	return docStyle.Render(m.list.View())
}

func (m model) headerView() string {
	createdAt := "Untitled"
	if m.viewingJournal != nil {
		createdAt = m.viewingJournal.CreatedAt.Format(timeFormat)
	}
	title := titleStyle.Render(createdAt)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) renderJournalView() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
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
	case JournalView:
		s = m.renderJournalView()
	}

	return s
}

func Root(ctx context.Context, js ports.JournalService) {
	p := tea.NewProgram(initialModel(ctx, js), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
