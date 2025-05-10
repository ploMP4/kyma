package transitions

import tea "github.com/charmbracelet/bubbletea"

type noTransition struct{}

func newNoTransition(_ int) noTransition {
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

func (t noTransition) Name() string {
	return "none"
}

func (t noTransition) Opposite() Transition {
	return t
}

func (t noTransition) Direction() direction {
	// don't care, no anim
	return Forwards
}
