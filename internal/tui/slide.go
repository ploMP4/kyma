package tui

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/goccy/go-yaml"
)

type propertyConfig struct {
	Style struct {
		Border      string `yaml:"border"`
		BorderColor string `yaml:"border_color"`
	} `yaml:"style"`
	Transition string `yaml:"transition"`
}

func getBorder(border string) lipgloss.Border {
	switch border {
	case "rounded":
		return lipgloss.RoundedBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "hidden":
		return lipgloss.HiddenBorder()
	case "block":
		return lipgloss.BlockBorder()
	case "innerHalfBlock":
		return lipgloss.InnerHalfBlockBorder()
	case "outerHalfBlock":
		return lipgloss.OuterHalfBlockBorder()
	default:
		return lipgloss.NormalBorder()
	}
}

type Properties struct {
	Style      StyleConfig
	Transition Transition
}

func NewProperties(properties string) (Properties, error) {
	if properties == "" {
		return Properties{Transition: getTransition("default", fps)}, nil
	}

	var p propertyConfig
	if err := yaml.Unmarshal([]byte(properties), &p); err != nil {
		return Properties{}, err
	}

	return Properties{
		Style: StyleConfig{
			Border:      getBorder(p.Style.Border),
			BorderColor: p.Style.BorderColor,
		},
		Transition: getTransition(p.Transition, fps),
	}, nil
}

type Slide struct {
	Data       string
	Prev       *Slide
	Next       *Slide
	Style      lipgloss.Style
	Properties Properties
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

	if s.Properties.Transition != nil && s.Properties.Transition.Animating() {
		b.WriteString(s.Properties.Transition.View(s.Prev.View(), s.Style.Render(out)))
	} else {
		b.WriteString(s.Style.Render(out))
	}
	return b.String()
}
