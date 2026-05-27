package widgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	dot   *StatusDot
	label *widget.Label
	root  *fyne.Container
}

func NewStatusBar() *StatusBar {
	dot := NewStatusDot()
	label := widget.NewLabel("● MAIN_MENU  │  ⚔️ —  │  ▶ Play  │  Раунд 0")
	root := container.NewHBox(dot.CanvasObject(), label)
	return &StatusBar{dot: dot, label: label, root: root}
}

func (s *StatusBar) Update(state, mode, action string, round int) {
	if mode == "" {
		mode = "—"
	}
	s.dot.SetActive(action == "pause")
	s.label.SetText(fmt.Sprintf("● %s  │  ⚔️ %s  │  %s  │  Раунд %d", state, mode, action, round))
}

func (s *StatusBar) Widget() fyne.CanvasObject { return s.root }
