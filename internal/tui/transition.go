package tui

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"

	"github.com/ploMP4/orama/internal/tui/messages"
)

type Transition interface {
	Start(width, height int) Transition
	Animating() bool
	Update() (Transition, tea.Cmd)
	View(prev, next string) string
}

type verticalSlideTransition struct {
	height    int
	fps       int
	spring    harmonica.Spring
	y         float64
	yVel      float64
	animating bool
}

func NewVerticalSlideTransition(fps int) verticalSlideTransition {
	const frequency = 7.0
	const damping = 0.9

	return verticalSlideTransition{
		fps:    fps,
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}
}

func (t verticalSlideTransition) Start(_, height int) Transition {
	t.height = height
	t.animating = true
	t.y = 0
	t.yVel = 0
	return t
}

func (t verticalSlideTransition) Animating() bool {
	return t.animating
}

func (t verticalSlideTransition) Update() (Transition, tea.Cmd) {
	targetY := float64(t.height)

	t.y, t.yVel = t.spring.Update(t.y, t.yVel, targetY)

	if t.y >= targetY {
		t.animating = false
		return t, nil
	}

	return t, messages.Animate(time.Duration(t.fps))
}

func (t verticalSlideTransition) View(prev, next string) string {
	var s strings.Builder

	y := int(math.Round(t.y))

	s.WriteString(prev)

	if t.animating {
		lines := strings.Split(next, "\n")
		if y > len(lines) {
			y = len(lines)
		}
		s.WriteString("\n" + strings.Join(lines[:y], "\n"))
	}

	return s.String()
}
