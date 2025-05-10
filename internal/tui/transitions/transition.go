package transitions

import tea "github.com/charmbracelet/bubbletea"

type Transition interface {
	Start(width, height int) Transition
	Animating() bool
	Update() (Transition, tea.Cmd)
	View(prev, next string) string
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
