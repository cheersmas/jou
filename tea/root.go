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

type errMsg error

type View string

const (
	Menu    View = "Menu"
	Add     View = "Add"
	List    View = "View"
	Edit    View = "Edit"
	Confirm View = "Confirm"

	// formattings
	timeFormat = "2 Jan, 2006"
)

type model struct {
	// meta
	context context.Context
	// service
	service ports.JournalService
	// main screen
	options        []View
	currentView    View
	cursorAtOption int

	// journals
	journals        []domains.Journal
	cursorAtJournal int
	selectedJournal domains.Journal

	// writing
	recentlySavedId int
	textarea        textarea.Model
	textError       error
}

func initialModel(ctx context.Context, js ports.JournalService) model {
	// handle text area
	ti := textarea.New()

	return model{
		context:         ctx,
		options:         []View{Add, List, Edit},
		currentView:     Menu,
		service:         js,
		textarea:        ti,
		textError:       nil,
		recentlySavedId: -1,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) setJournals() model {
	ctx := context.Background()
	journals, err := m.service.ListAll(ctx)
	m.journals = journals

	if err != nil {
		log.Fatalf("Something went wrong fetching the results: %v", err)
		panic(err)
	}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// quitting
		case "esc":
			switch m.currentView {
			case Add:
				if m.textarea.Focused() {
					m.textarea.Blur()
				}
			case Confirm:
				m.currentView = Add
			}
		case "ctrl+c", "q":
			switch m.currentView {
			case Add:
				m.currentView = Confirm
				return m, nil
			default:
				return m, tea.Quit
			}
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			switch m.currentView {
			case Menu:
				if m.cursorAtOption > 0 {
					m.cursorAtOption--
				}
			case List, Edit:
				if m.cursorAtJournal > 0 {
					m.cursorAtJournal--
				}
			}
		// The "down" and "j" keys move the cursor down
		case "down", "j":
			switch m.currentView {
			case Menu:
				if m.cursorAtOption < len(m.options)-1 {
					m.cursorAtOption++
				}
			case List, Edit:
				if m.cursorAtJournal < len(m.journals)-1 {
					m.cursorAtJournal++
				}
			}
		case "enter":
			switch m.currentView {
			case Menu:
				m.currentView = m.options[m.cursorAtOption]
				if m.currentView == List || m.currentView == Edit {
					m = m.setJournals()
				}
			case List, Edit:
				m.selectedJournal = m.journals[m.cursorAtJournal]
				if m.currentView == Edit {
					m.textarea.SetValue(m.selectedJournal.Content)
					m.currentView = Add
				} else {
					// TODO: handle view only mode in a new window here
				}
			case Confirm:
				m.currentView = Menu
			}

		case "ctrl+s":
			switch m.currentView {
			case Add:
				textAreaValue := m.textarea.Value()
				journal := domains.Journal{
					Content: textAreaValue,
				}

				// draft flow
				if m.recentlySavedId == -1 {
					id, err := m.service.Create(m.context, journal)
					if err != nil {
						m.textError = err
					}
					m.recentlySavedId = id
				} else if m.recentlySavedId > -1 { // non draft flow
					id, err := m.service.Update(m.context, m.recentlySavedId, textAreaValue)
					if err != nil {
						m.textError = err
					}
					m.recentlySavedId = id
				}
			}
		case "backspace":
			switch m.currentView {
			case List, Edit:
				m.currentView = Menu
			}

		default:
			switch m.currentView {
			case Add:
				if !m.textarea.Focused() {
					cmd = m.textarea.Focus()
					cmds = append(cmds, cmd)
				}
			}
		}

	case errMsg:
		m.textError = msg
		return m, nil
	}

	if m.currentView == Add {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) showMenuView() string {
	s := "What would like to do?\n\n"
	for i, option := range m.options {
		cursor := "[ ]"
		if m.cursorAtOption == i {
			cursor = "[>]"
		}

		s += fmt.Sprintf("%s %s\n", cursor, option)
	}
	return s
}

func (m model) showConfirmationView() string {
	s := "Unsaved changes can get lost.\n\n Press: \n<esc> to cancel and go back \nor\n <enter> to discard"
	return s
}

func (m model) showAddJournalView() string {
	// we came from the list screen
	createdAtDate := time.Now().Format(timeFormat)
	if m.selectedJournal.Id > 0 {
		createdAtDate = m.selectedJournal.CreatedAt.Format(timeFormat)
	}
	s := fmt.Sprintf("Journal entry %s\n\n", createdAtDate)
	s += fmt.Sprintf(
		"%s",
		m.textarea.View(),
	)

	// TODO: move this to the footer of the recently saved
	if m.recentlySavedId > 0 {
		s += fmt.Sprintf("\nSaved %d\n", m.recentlySavedId)
	}
	return s
}

func (m model) showListJournalsView() string {
	s := "All Journals\n\n"

	for i, journal := range m.journals {
		cursor := "[ ]"
		if i == m.cursorAtJournal {
			cursor = "[>]"
		}
		formatedDate := journal.CreatedAt.Format(timeFormat)
		s += fmt.Sprintf("%s: %s | %s \n", cursor, formatedDate, journal.Content)
	}

	s += "\n"
	return s
}

func (m model) View() string {

	s := ""
	// show menu when nothing is selected
	switch m.currentView {
	case Menu:
		s += m.showMenuView()
	case Add:
		s += m.showAddJournalView()
	case List, Edit:
		s += m.showListJournalsView()
	case Confirm:
		s += m.showConfirmationView()
	}

	if m.currentView != Confirm {
		s += "\n(Press ctrl+c or q to quit)\n"
	}

	return s
}

func Root(ctx context.Context, js ports.JournalService) {
	p := tea.NewProgram(initialModel(ctx, js))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
