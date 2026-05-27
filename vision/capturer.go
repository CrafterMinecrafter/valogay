package vision
import "image"
type Capturer interface { CaptureRect(rect image.Rectangle)(image.Image,error); DirtyRects() []image.Rectangle; Close() error }
