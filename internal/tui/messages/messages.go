package messages

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type FrameMsg time.Time

func Animate(fps time.Duration) tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return FrameMsg(t)
	})
}

func Wait(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return nil
	}
}
