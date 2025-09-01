package app

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cheersmas/jou/app/constants"
	"github.com/cheersmas/jou/app/input"
	"github.com/cheersmas/jou/app/navigation"
	"github.com/cheersmas/jou/app/styles"
	"github.com/cheersmas/jou/app/views"
	"github.com/cheersmas/jou/ports"
)

type App struct {
	state        *navigation.AppState
	router       *navigation.Router
	inputHandler *input.InputHandler
	views        map[constants.View]views.View
}

func NewApp(ctx context.Context, service ports.JournalService) *App {
	state := navigation.NewAppState(ctx, service)
	router := navigation.NewRouter(state)
	inputHandler := input.NewInputHandler(state, router)

	// Initialize UI components
	ti := textarea.New()
	ti.Placeholder = "Write your journal entry here..."

	vp := viewport.New(30, 5)
	vp.SetContent(`init`)

	items := []list.Item{}
	li := list.New(items, list.NewDefaultDelegate(), 0, 0)

	li.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "back to menu"),
			),
		}
	}

	li.DisableQuitKeybindings()

	state.Textarea = ti
	state.Viewport = vp
	state.List = li

	// Initialize views
	viewMap := map[constants.View]views.View{
		constants.MenuView:    views.MenuView{},
		constants.AddView:     views.AddView{},
		constants.EditView:    views.ListView{},
		constants.ListView:    views.ListView{},
		constants.JournalView: views.JournalView{},
		constants.ConfirmView: views.ConfirmView{},
	}

	return &App{
		state:        state,
		router:       router,
		inputHandler: inputHandler,
		views:        viewMap,
	}
}

func (a App) Init() tea.Cmd {
	return textarea.Blink
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmd = a.inputHandler.HandleKeyMsg(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		a.handleWindowSize(msg)
	case error:
		a.state.LastError = msg
		return a, nil
	}

	// Update current view
	if view, exists := a.views[a.state.CurrentView]; exists {
		cmd = view.Update(a.state, msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if view, exists := a.views[a.state.CurrentView]; exists {
		return view.Render(a.state)
	}
	return "Unknown view"
}

func (a *App) handleWindowSize(msg tea.WindowSizeMsg) {
	// Handle window size changes
	h, v := styles.DocStyle.GetFrameSize()
	a.state.List.SetSize(msg.Width-h, msg.Height-v)

	// Update textarea size
	a.state.Textarea.SetHeight(msg.Height - 10) // Adjust as needed
	a.state.Textarea.SetWidth(msg.Width)

	if !a.state.Ready {
		a.state.Viewport = viewport.New(msg.Width, msg.Height-10)
		a.state.Ready = true
	} else {
		a.state.Viewport.Width = msg.Width
		a.state.Viewport.Height = msg.Height - 10
	}
}

func Root(ctx context.Context, js ports.JournalService) {
	app := NewApp(ctx, js)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
