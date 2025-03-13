package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/ploMP4/orama/internal/tui/messages"
)

type keyMap struct {
	Quit key.Binding
	Next key.Binding
	Prev key.Binding
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
	Next: key.NewBinding(
		key.WithKeys("right", "l", " "),
		key.WithHelp(">, l, <SPC>", "next"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<, h", "previous"),
	),
}

const (
	fps       = 60
	frequency = 7.0
	damping   = 0.9
)

func style(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9999CC")). // Blueish
		Width(width - 4).
		Height(height - 2)
}

type Slide struct {
	Data       string
	Prev       *Slide
	Next       *Slide
	Style      lipgloss.Style
	Transition Transition
}

func (s Slide) View() string {
	var b strings.Builder

	out, err := glamour.Render(s.Data, "dark")
	if err != nil {
		b.WriteString("\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Render("Error: "+err.Error()))
		return b.String()
	}

	if s.Transition != nil && s.Transition.Animating() {
		b.WriteString(s.Transition.View(s.Prev.View(), s.Style.Render(out)))
	} else {
		b.WriteString(s.Style.Render(out))
	}
	return b.String()
}

type model struct {
	width  int
	height int

	slide *Slide
	keys  keyMap
	help  help.Model
}

func New(rootSlide *Slide) model {
	return model{
		slide: rootSlide,
		keys:  keys,
		help:  help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		slide := m.slide
		for slide != nil {
			slide.Style = style(m.width, m.height)
			slide = slide.Next
		}
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		} else if key.Matches(msg, m.keys.Next) {
			if m.slide.Next == nil {
				return m, nil
			}
			m.slide = m.slide.Next
			m.slide.Transition = m.slide.Transition.Start(m.width, m.height)
			return m, messages.Animate(fps)
		} else if key.Matches(msg, m.keys.Prev) {
			if m.slide.Prev == nil {
				return m, nil
			}
			m.slide = m.slide.Prev
			return m, nil
		}
	case messages.FrameMsg:
		transition, cmd := m.slide.Transition.Update()
		m.slide.Transition = transition
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		m.slide.View(),
	)
}
