//go:build linux

package vision

import (
	"errors"
	"image"
)

type Capturer interface {
	CaptureRect(rect image.Rectangle) (image.Image, error)
	DirtyRects() []image.Rectangle
	Close() error
}

type X11Capturer struct{}

func NewDXGICapturer(_ int) (*X11Capturer, error) {
	return &X11Capturer{}, nil
}

func (x *X11Capturer) CaptureRect(rect image.Rectangle) (image.Image, error) {
	if rect.Empty() {
		return nil, errors.New("empty rect")
	}
	return image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy())), nil
}

func (x *X11Capturer) DirtyRects() []image.Rectangle { return nil }

func (x *X11Capturer) Close() error { return nil }
