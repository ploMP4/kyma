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

type noTransition struct{}

func NewNoTransition() noTransition {
	return noTransition{}
}

func (t noTransition) Start(width int, height int) Transition {
	return t
}

func (t noTransition) Animating() bool {
	return false
}

func (t noTransition) Update() (Transition, tea.Cmd) {
	return t, nil
}

func (notransition noTransition) View(prev string, next string) string {
	return ""
}

type verticalSlideUpTransition struct {
	height    int
	fps       int
	spring    harmonica.Spring
	y         float64
	yVel      float64
	animating bool
}

func NewVerticalSlideUpTransition(fps int) verticalSlideUpTransition {
	const frequency = 7.0
	const damping = 0.8

	return verticalSlideUpTransition{
		fps:    fps,
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}
}

func (t verticalSlideUpTransition) Start(_, height int) Transition {
	t.height = height
	t.animating = true
	t.y = 0
	t.yVel = 0
	return t
}

func (t verticalSlideUpTransition) Animating() bool {
	return t.animating
}

func (t verticalSlideUpTransition) Update() (Transition, tea.Cmd) {
	targetY := float64(t.height)

	t.y, t.yVel = t.spring.Update(t.y, t.yVel, targetY)

	if t.y >= targetY {
		t.animating = false
		return t, nil
	}

	return t, messages.Animate(time.Duration(t.fps))
}

func (t verticalSlideUpTransition) View(prev, next string) string {
	var s strings.Builder

	y := int(math.Round(t.y))

	s.WriteString(prev)

	lines := strings.Split(next, "\n")
	if y > len(lines) {
		y = len(lines)
	}
	s.WriteString("\n" + strings.Join(lines[:y], "\n"))

	return s.String()
}

type verticalSlideDownTransition struct {
	height    int
	fps       int
	spring    harmonica.Spring
	y         float64
	yVel      float64
	animating bool
}

func NewVerticalSlideDownTransition(fps int) verticalSlideDownTransition {
	const frequency = 7.0
	const damping = 0.8

	return verticalSlideDownTransition{
		fps:    fps,
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}
}

func (t verticalSlideDownTransition) Start(_, height int) Transition {
	t.height = height
	t.animating = true
	t.y = 0
	t.yVel = 0
	return t
}

func (t verticalSlideDownTransition) Animating() bool {
	return t.animating
}

func (t verticalSlideDownTransition) Update() (Transition, tea.Cmd) {
	targetY := float64(t.height)

	t.y, t.yVel = t.spring.Update(t.y, t.yVel, targetY)

	if t.y >= targetY {
		t.animating = false
		return t, nil
	}

	return t, messages.Animate(time.Duration(t.fps))
}

func (t verticalSlideDownTransition) View(prev, next string) string {
	var s strings.Builder

	y := int(math.Round(t.y))

	nextLines := strings.Split(next, "\n")
	if y > len(nextLines) {
		y = len(nextLines)
	}
	s.WriteString(strings.Join(nextLines[len(nextLines)-y:], "\n"))

	prevLines := strings.Split(prev, "\n")
	if y > len(prevLines) {
		y = len(prevLines)
	}
	s.WriteString("\n" + strings.Join(prevLines[:len(prevLines)-y], "\n"))

	return s.String()
}
