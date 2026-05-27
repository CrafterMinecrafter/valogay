//go:build !windows

package vision

import "image"

type Capturer interface {
	CaptureRect(rect image.Rectangle) (image.Image, error)
	DirtyRects() []image.Rectangle
	Close() error
}

type DXGICapturer struct{}

func NewDXGICapturer(_ int) (*DXGICapturer, error) { return &DXGICapturer{}, nil }
func (d *DXGICapturer) CaptureRect(rect image.Rectangle) (image.Image, error) {
	return image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy())), nil
}
func (d *DXGICapturer) DirtyRects() []image.Rectangle { return nil }
func (d *DXGICapturer) Close() error                  { return nil }
