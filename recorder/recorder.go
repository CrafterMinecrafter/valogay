//go:build windows

package recorder

import (
	"image"
	"sync"
	"time"
)

type Frame struct {
	Pixels    []uint8
	Width     int
	Height    int
	Timestamp int64
}

func New(intervalSec int, outputDir string) *Recorder {
	return &Recorder{
		interval:  time.Duration(intervalSec) * time.Second,
		stopCh:    make(chan struct{}),
		outputDir: outputDir,
	}
}

const (
	maxFrames   = 120
	scaleFactor = 6
)

type Recorder struct {
	interval  time.Duration
	stopCh    chan struct{}
	outputDir string

	mu      sync.RWMutex
	frames  []Frame
	running bool
}

func (r *Recorder) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.running {
		return
	}
	r.running = true
	r.stopCh = make(chan struct{})
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				img, err := captureScreen()
				if err != nil {
					continue
				}
				saveFrame(img, r.outputDir)
				frame := toRedGray(img)
				r.mu.Lock()
				r.frames = append(r.frames, frame)
				if len(r.frames) > maxFrames {
					r.frames = r.frames[1:]
				}
				r.mu.Unlock()
			case <-r.stopCh:
				r.mu.Lock()
				r.running = false
				r.mu.Unlock()
				return
			}
		}
	}()
}

func (r *Recorder) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.running {
		return
	}
	close(r.stopCh)
}

func (r *Recorder) Frames() []Frame {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Frame, len(r.frames))
	copy(out, r.frames)
	return out
}

func (r *Recorder) LastFrame() *Frame {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.frames) == 0 {
		return nil
	}
	f := r.frames[len(r.frames)-1]
	return &f
}

func toRedGray(img *image.RGBA) Frame {
	b := img.Bounds()
	w := b.Dx() / scaleFactor
	h := b.Dy() / scaleFactor
	pixels := make([]uint8, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx := x * scaleFactor
			sy := y * scaleFactor
			off := (sy-b.Min.Y)*img.Stride + (sx-b.Min.X)*4
			pixels[y*w+x] = img.Pix[off]
		}
	}
	return Frame{Pixels: pixels, Width: w, Height: h, Timestamp: time.Now().UnixMilli()}
}

type Rect struct {
	Left, Top, Right, Bottom int32
}

func (r Rect) ToImage() image.Rectangle {
	return image.Rect(int(r.Left), int(r.Top), int(r.Right), int(r.Bottom))
}