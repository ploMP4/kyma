package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

const fps = 60

func style(width, height int, extra StyleConfig) lipgloss.Style {
	borderColor := "#9999CC" // Blueish
	if extra.BorderColor != "" {
		borderColor = extra.BorderColor
	}

	return extra.Layout.
		Border(extra.Border).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(width - 4).
		Height(height - 2)
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
			slide.Style = style(m.width, m.height, slide.Properties.Style)
			slide = slide.Next
		}
		return m, nil
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		} else if key.Matches(msg, m.keys.Next) {
			if m.slide.Next == nil || m.slide.Properties.Transition.Animating() {
				return m, nil
			}
			m.slide = m.slide.Next
			m.slide.Properties.Transition = m.slide.Properties.Transition.Start(m.width, m.height)
			return m, messages.Animate(fps)
		} else if key.Matches(msg, m.keys.Prev) {
			if m.slide.Prev == nil || m.slide.Properties.Transition.Animating() {
				return m, nil
			}
			m.slide = m.slide.Prev
			return m, messages.Animate(fps)
		}
	case messages.FrameMsg:
		slide, cmd := m.slide.Update()
		m.slide = slide
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
