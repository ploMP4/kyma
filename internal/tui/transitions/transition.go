package transitions

import tea "github.com/charmbracelet/bubbletea"

type direction byte

const (
	Forwards  direction = 0
	Backwards direction = 1
)

type Transition interface {
	Start(width, height int) Transition
	Animating() bool
	Update() (Transition, tea.Cmd)
	View(prev, next string) string
	Opposite() Transition
	Name() string
}

func Get(name string, fps int) Transition {
	switch name {
	case "slideUp":
		return newSlideUp(fps)
	case "slideDown":
		return newSlideDown(fps)
	case "swipeLeft":
		return newSwipeLeft(fps)
	case "swipeRight":
		return newSwipeRight(fps)
	case "flip":
		return newFlipRight(fps)
	default:
		return newNoTransition(fps)
	}
}
