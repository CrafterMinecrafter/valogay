//go:build windows

package vision

import (
	"errors"
	"image"
	"sync"
)

type Capturer interface {
	CaptureRect(rect image.Rectangle) (image.Image, error)
	DirtyRects() []image.Rectangle
	Close() error
}

type DXGICapturer struct {
	bufPool   sync.Pool
	lastDirty []image.Rectangle
}

func NewDXGICapturer(_ int) (*DXGICapturer, error) {
	c := &DXGICapturer{}
	c.bufPool.New = func() any { return image.NewRGBA(image.Rect(0, 0, 1, 1)) }
	return c, nil
}

func (d *DXGICapturer) CaptureRect(rect image.Rectangle) (image.Image, error) {
	if rect.Empty() {
		return nil, errors.New("empty rect")
	}
	return image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy())), nil
}
func (d *DXGICapturer) DirtyRects() []image.Rectangle { return d.lastDirty }
func (d *DXGICapturer) Close() error                  { return nil }
