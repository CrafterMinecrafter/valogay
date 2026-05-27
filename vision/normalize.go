package vision

import (
	"image"
	"image/color"
)

func Normalize(img image.Image) image.Image {
	b := img.Bounds()
	out := image.NewGray(image.Rect(0, 0, 8, 8))
	w, h := b.Dx(), b.Dy()
	if w == 0 || h == 0 {
		return out
	}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			sx := b.Min.X + x*w/8
			sy := b.Min.Y + y*h/8
			r, g, bl, _ := img.At(sx, sy).RGBA()
			gray := uint8(((r + g + bl) / 3) >> 8)
			out.Set(x, y, color.Gray{Y: gray})
		}
	}
	return out
}
