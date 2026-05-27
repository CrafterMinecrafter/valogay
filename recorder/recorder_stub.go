//go:build !windows

package recorder

type Frame struct {
	Pixels    []uint8
	Width     int
	Height    int
	Timestamp int64
}

type Recorder struct{}

func New(intervalSec int, outputDir string) *Recorder {
	return &Recorder{}
}

func (r *Recorder) Start()            {}
func (r *Recorder) Stop()             {}
func (r *Recorder) Frames() []Frame   { return nil }
func (r *Recorder) LastFrame() *Frame { return nil }
