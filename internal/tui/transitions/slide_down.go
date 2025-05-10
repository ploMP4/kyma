package transitions

import (
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"

	"github.com/ploMP4/kyma/internal/tui/messages"
)

type slideDown struct {
	height    int
	fps       int
	spring    harmonica.Spring
	y         float64
	yVel      float64
	animating bool
}

func newSlideDown(fps int) slideDown {
	const frequency = 7.0
	const damping = 0.8

	return slideDown{
		fps:    fps,
		spring: harmonica.NewSpring(harmonica.FPS(fps), frequency, damping),
	}
}

func (t slideDown) Start(_, height int) Transition {
	t.height = height
	t.animating = true
	t.y = 0
	t.yVel = 0
	return t
}

func (t slideDown) Animating() bool {
	return t.animating
}

func (t slideDown) Update() (Transition, tea.Cmd) {
	targetY := float64(t.height)

	t.y, t.yVel = t.spring.Update(t.y, t.yVel, targetY)

	if t.y >= targetY {
		t.animating = false
		return t, nil
	}

	return t, messages.Animate(time.Duration(t.fps))
}

func (t slideDown) View(prev, next string) string {
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

func (t slideDown) Name() string {
	return "slideDown"
}
