package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/goccy/go-yaml"

	"github.com/ploMP4/kyma/internal/tui/transitions"
)

type Slide struct {
	Data       string
	Prev       *Slide
	Next       *Slide
	Style      lipgloss.Style
	Properties Properties

	preRenderedFrame string
}

func (s *Slide) Update() (*Slide, tea.Cmd) {
	transition, cmd := s.Properties.Transition.Update()
	s.Properties.Transition = transition
	s.preRenderedFrame = s.view()
	if cmd == nil {
		s.preRenderedFrame = ""
	}
	return s, cmd
}

func (s Slide) View() string {
	if s.preRenderedFrame == "" {
		return s.view()
	}
	return s.preRenderedFrame
}

func (s Slide) view() string {
	var b strings.Builder

	out, err := glamour.Render(s.Data, "dark")
	if err != nil {
		b.WriteString("\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Render("Error: "+err.Error()))
		return b.String()
	}

	if s.Properties.Transition != nil && s.Properties.Transition.Animating() {
		b.WriteString(s.Properties.Transition.View(s.Prev.View(), s.Style.Render(out)))
	} else {
		b.WriteString(s.Style.Render(out))
	}
	return b.String()
}

type Properties struct {
	Style      StyleConfig            `yaml:"style"`
	Transition transitions.Transition `yaml:"transition"`
}

func (p *Properties) UnmarshalYAML(bytes []byte) error {
	aux := struct {
		Style      StyleConfig `yaml:"style"`
		Transition string      `yaml:"transition"`
	}{}

	if err := yaml.Unmarshal(bytes, &aux); err != nil {
		return err
	}
	p.Transition = transitions.Get(aux.Transition, fps)
	p.Style = aux.Style

	return nil
}

func NewProperties(properties string) (Properties, error) {
	if properties == "" {
		return Properties{Transition: transitions.Get("default", fps)}, nil
	}

	var p Properties
	if err := yaml.Unmarshal([]byte(properties), &p); err != nil {
		return Properties{}, err
	}

	return p, nil
}
