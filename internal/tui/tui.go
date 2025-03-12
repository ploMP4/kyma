package tui

import (
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/harmonica"
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

type model struct {
	width  int
	height int

	slides       []string
	currentSlide int
	keys         keyMap
	help         help.Model

	spring    harmonica.Spring
	y         float64
	yVel      float64
	animating bool
}

func New(slides []string) model {
	return model{
		slides:       slides,
		currentSlide: 0,
		keys:         keys,
		help:         help.New(),
		spring:       harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
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
		} else if key.Matches(msg, m.keys.Next) {
			if m.currentSlide == len(m.slides)-1 || m.animating {
				return m, nil
			}
			m.animating = true
			m.y = 0
			m.yVel = 0
			return m, messages.Animate(fps)
		} else if key.Matches(msg, m.keys.Prev) {
			if m.currentSlide == 0 || m.animating {
				return m, nil
			}
			m.currentSlide--
			return m, nil
		}
	case messages.FrameMsg:
		targetY := float64(m.height)

		m.y, m.yVel = m.spring.Update(m.y, m.yVel, targetY)

		if m.y >= targetY {
			m.animating = false
			m.currentSlide++
			return m, nil
		}

		return m, messages.Animate(fps)
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	y := int(math.Round(m.y))

	layout := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#9999CC")). // Blueish
		Width(m.width - 4).
		Height(m.height - 2)

	out, err := glamour.Render(m.slides[m.currentSlide], "dark")
	if err != nil {
		s.WriteString("\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Render("Error: "+err.Error()))
		return s.String()
	}

	s.WriteString(layout.Render(out))

	if m.animating {
		outNext, err := glamour.Render(m.slides[m.currentSlide+1], "dark")
		if err != nil {
			s.WriteString("\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")). // Red
				Render("Error: "+err.Error()))
			return s.String()
		}

		lines := strings.Split(layout.Render(outNext), "\n")
		if y > len(lines) {
			y = len(lines)
		}
		s.WriteString("\n" + strings.Join(lines[:y], "\n"))
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		s.String(),
	)
}
