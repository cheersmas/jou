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

type Option string

const (
	Menu Option = "Menu"
	Add  Option = "Add"
	View Option = "View"
)

type model struct {
	// service
	service ports.JournalService
	// main screen
	options        []Option
	selectedOption Option
	cursorAtOption int

	// journals
	journals        []domains.Journal
	selectedJournal int

	// writing
	textarea  textarea.Model
	textError error
}

func initialModel(ctx context.Context, js ports.JournalService) model {
	// handle text area
	ti := textarea.New()

	return model{
		options:        []Option{Add, View},
		selectedOption: Menu,
		service:        js,
		textarea:       ti,
		textError:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmds []tea.Cmd
// 	var cmd tea.Cmd

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyEsc:
// 			if m.textarea.Focused() {
// 				m.textarea.Blur()
// 			}
// 		case tea.KeyCtrlC:
// 			return m, tea.Quit

// 		default:
// 			if !m.textarea.Focused() {
// 				cmd = m.textarea.Focus()
// 				cmds = append(cmds, cmd)
// 			}
// 		}

// 	case errMsg:
// 		m.textError = msg
// 		return m, nil
// 	}

// 	m.textarea, cmd = m.textarea.Update(msg)
// 	cmds = append(cmds, cmd)
// 	return m, tea.Batch(cmds...)
// }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.selectedOption == Add {
				if m.textarea.Focused() {
					m.textarea.Blur()
				}
			}
		case "ctrl+c":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.selectedOption == Menu {
				if m.cursorAtOption > 0 {
					m.cursorAtOption--
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.selectedOption == Menu {
				if m.cursorAtOption < len(m.options)-1 {
					m.cursorAtOption++
				}
			}
		case "enter":
			if m.selectedOption == Menu {
				m.selectedOption = m.options[m.cursorAtOption]
			}

		case "backspace":
			if m.selectedOption != Add {
				m.selectedOption = Menu
			}

		default:
			if m.selectedOption == Add && !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case errMsg:
		m.textError = msg
		return m, nil
	}

	if m.selectedOption == Add {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) showMenu() string {
	s := "What would like to do?\n\n"
	for i, option := range m.options {
		cursor := " "
		if m.cursorAtOption == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, option)
	}
	return s
}
func (m model) showTextArea() string {
	createdAtDate := time.Now().Format("2006-01-02")
	s := fmt.Sprintf("Journal entry %s\n\n", createdAtDate)
	s += fmt.Sprintf(
		"%s",
		m.textarea.View(),
	)

	return s
}
func (m model) showJournals() string {
	ctx := context.Background()
	journals, err := m.service.ListAll(ctx)
	if err != nil {
		log.Fatalf("Something went wrong fetching the results: %v", err)
		panic(err)
	}

	s := "All Journals\n\n"

	for i, journal := range journals {
		formatedDate := journal.CreatedAt.Format("2006-01-02 15:04")
		s += fmt.Sprintf("%d: %s | %s \n", i+1, formatedDate, journal.Content)
	}

	s += "\n"
	return s
}

func (m model) View() string {

	s := ""
	// show menu when nothing is selected
	switch m.selectedOption {
	case Menu:
		s += m.showMenu()
	case Add:
		s += m.showTextArea()
	case View:
		s += m.showJournals()
	}

	s += "\n(Press ctrl+c or q to quit)\n"
	return s
}

func Root(ctx context.Context, js ports.JournalService) {
	p := tea.NewProgram(initialModel(ctx, js))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
