package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q, esc, ctrl+c", "quit"),
	),
}

type model struct {
	width  int
	height int

	slides       []string
	currentSlide uint
	keys         keyMap
	help         help.Model
}

func New(slides []string) model {
	return model{
		slides:       slides,
		currentSlide: 0,
		keys:         keys,
		help:         help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	layout := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9999CC")). // Blueish
		PaddingLeft(2).
		PaddingRight(2)

	out, err := glamour.Render(m.slides[m.currentSlide], "dark")
	if err != nil {
		s.WriteString("\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Render("Error: "+err.Error()))
		return s.String()
	}

	s.WriteString(out)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		layout.Render(s.String()),
	)
}
