package widgets

import "image"

type RegionEditor struct {
	screenshot  image.Image
	selection   image.Rectangle
	zoom        float32
	dragging    bool
	dragStart   image.Point
	liveTest    bool
	liveCapture image.Image
	OnSave      func(rect image.Rectangle, threshold int, refPath string)
}
