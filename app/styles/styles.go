package styles

import "github.com/charmbracelet/lipgloss"

var (
	DocStyle = lipgloss.NewStyle().Margin(1, 2)

	RoundedBorderStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	TitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	InfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return TitleStyle.BorderStyle(b)
	}()

	ContainerStyle = lipgloss.NewStyle().Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)
