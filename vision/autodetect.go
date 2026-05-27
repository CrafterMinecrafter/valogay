package vision

import (
	"image"
	"image/color"
)

func AutoCropToContent(_ image.Image, roughRect image.Rectangle, _ int) image.Rectangle {
	return roughRect
}
func FindByDominantColor(_ image.Image, roughRect image.Rectangle, _ color.RGBA, _ int) image.Rectangle {
	return roughRect
}
