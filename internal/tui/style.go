package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/goccy/go-yaml"
)

type StyleConfig struct {
	// Layout      lipgloss.Style  `yaml:"layout"`
	Border      lipgloss.Border `yaml:"border"`
	BorderColor string          `yaml:"border_color"`
}

func (s *StyleConfig) UnmarshalYAML(bytes []byte) error {
	aux := struct {
		// Layout      string `yaml:"layout"`
		Border      string `yaml:"border"`
		BorderColor string `yaml:"border_color"`
	}{}

	if err := yaml.Unmarshal(bytes, &aux); err != nil {
		return err
	}
	// s.Layout = getLayout(aux.Layout)
	s.Border = getBorder(aux.Border)
	s.BorderColor = aux.BorderColor

	return nil
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

// func getLayout(layout string) lipgloss.Style {
// 	switch layout {
// 	case "center":
// 		return lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center)
// 	case "left":
// 		return lipgloss.NewStyle().Align(lipgloss.Left, lipgloss.Left)
// 	case "right":
// 		return lipgloss.NewStyle().Align(lipgloss.Right, lipgloss.Center)
// 	default:
// 		return lipgloss.NewStyle()
// 	}
// }
