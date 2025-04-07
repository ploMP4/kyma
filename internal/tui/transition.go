package tui

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/reflow/wordwrap"

	"github.com/ploMP4/orama/internal/skip"
	"github.com/ploMP4/orama/internal/tui/messages"
)

type Transition interface {
	Start(width, height int) Transition
	Animating() bool
	Update() (Transition, tea.Cmd)
	View(prev, next string) string
}

func getTransition(name string, fps int) Transition {
	switch name {
	case "verticalUp":
		return NewVerticalSlideUpTransition(fps)
	case "verticalDown":
		return NewVerticalSlideDownTransition(fps)
	case "swipe":
		return NewHorizontalSlideLeftTransition(fps)
	case "flip":
		return NewPageFlipRightTransition(fps)
	default:
		return NewNoTransition(fps)
	}
}

type noTransition struct{}

func NewNoTransition(_ int) noTransition {
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

type horizontalSlideLeftTransition struct {
	width     int
	fps       int
	spring    harmonica.Spring
	x         float64
	xVel      float64
	animating bool
}

func NewHorizontalSlideLeftTransition(fps int) horizontalSlideLeftTransition {
	const frequency = 7.0
	const damping = 0.8

	return horizontalSlideLeftTransition{
		fps:    fps,
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}
}

func (t horizontalSlideLeftTransition) Start(width int, _ int) Transition {
	t.width = width
	t.animating = true
	t.x = 0
	t.xVel = 0
	return t
}

func (t horizontalSlideLeftTransition) Animating() bool {
	return t.animating
}

func (t horizontalSlideLeftTransition) Update() (Transition, tea.Cmd) {
	targetX := float64(t.width)

	t.x, t.xVel = t.spring.Update(t.x, t.xVel, targetX)

	if t.x >= targetX {
		t.animating = false
		return t, nil
	}

	return t, messages.Animate(time.Duration(t.fps))
}

func (t horizontalSlideLeftTransition) View(prev string, next string) string {
	var s strings.Builder

	x := int(math.Round(t.x))

	prevLines := strings.Split(prev, "\n")
	nextLines := strings.Split(next, "\n")

	// assert that slides are always equal size
	if len(nextLines) != len(prevLines) {
		panic("Slides of not equal height")
	}

	for i := range nextLines {
		prev := skip.String(prevLines[i], uint(x))
		next := truncate.String(nextLines[i], uint(x))
		s.WriteString(prev + " " + next)
		if i < len(nextLines)-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

type pageFlipRightTransition struct {
	width     int
	fps       int
	spring    harmonica.Spring
	x         float64
	xVel      float64
	animating bool
}

func NewPageFlipRightTransition(fps int) pageFlipRightTransition {
	const frequency = 7.0
	const damping = 0.8

	return pageFlipRightTransition{
		fps:    fps,
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}
}

func (t pageFlipRightTransition) Start(width int, _ int) Transition {
	t.width = width
	t.animating = true
	t.x = 0
	t.xVel = 0
	return t
}

func (t pageFlipRightTransition) Animating() bool {
	return t.animating
}

func (t pageFlipRightTransition) Update() (Transition, tea.Cmd) {
	targetX := float64(t.width)

	t.x, t.xVel = t.spring.Update(t.x, t.xVel, targetX)

	if t.x >= targetX {
		t.animating = false
		return t, nil
	}

	return t, messages.Animate(time.Duration(t.fps))
}

func (t pageFlipRightTransition) View(prev string, next string) string {
	var s strings.Builder

	x := int(math.Round(t.x))

	prevLines := strings.Split(prev, "\n")
	nextLines := strings.Split(next, "\n")

	// assert that slides are always equal size
	if len(nextLines) != len(prevLines) {
		panic("Slides of not equal height")
	}

	for i := range nextLines {
		wrappedPrev := strings.Split(wordwrap.String(prevLines[i], x), "\n")
		prev = wrappedPrev[len(wrappedPrev)-1]
		next := truncate.String(nextLines[i], uint(x))
		s.WriteString(next + " " + prev)
		if i < len(nextLines)-1 {
			s.WriteString("\n")
		}
	}

	return s.String()
}
