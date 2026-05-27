package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type StatusDot struct {
	Active bool
	dot    *canvas.Circle
}

func NewStatusDot() *StatusDot {
	d := &StatusDot{dot: canvas.NewCircle(color.NRGBA{R: 180, G: 180, B: 180, A: 255})}
	return d
}

func (s *StatusDot) SetActive(v bool) {
	s.Active = v
	if v {
		s.dot.FillColor = color.NRGBA{R: 30, G: 210, B: 90, A: 255}
	} else {
		s.dot.FillColor = color.NRGBA{R: 210, G: 70, B: 70, A: 255}
	}
	s.dot.Refresh()
}

func (s *StatusDot) CanvasObject() fyne.CanvasObject { return s.dot }
