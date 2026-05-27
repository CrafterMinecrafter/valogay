package widgets

import (
	"image"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RegionEditor struct {
	screenshot image.Image
	selection  image.Rectangle
	zoom       float32
	pan        fyne.Position
	dragging   bool
	dragStart  image.Point
	liveTest   bool

	canvasImage *canvas.Image
	selectionR  *canvas.Rectangle
	root        *fyne.Container

	OnSave func(rect image.Rectangle, threshold int, refPath string)
}

func NewRegionEditor(img image.Image) *RegionEditor {
	r := &RegionEditor{screenshot: img, zoom: 1.0}
	r.buildUI()
	return r
}

func (r *RegionEditor) buildUI() {
	r.canvasImage = canvas.NewImageFromImage(r.screenshot)
	r.canvasImage.FillMode = canvas.ImageFillContain
	r.selectionR = canvas.NewRectangle(color.NRGBA{R: 60, G: 130, B: 255, A: 80})
	r.selectionR.StrokeColor = color.NRGBA{R: 60, G: 130, B: 255, A: 255}
	r.selectionR.StrokeWidth = 2
	layer := container.NewWithoutLayout(r.canvasImage, r.selectionR)

	x := widget.NewEntry()
	y := widget.NewEntry()
	w := widget.NewEntry()
	h := widget.NewEntry()
	threshold := widget.NewSlider(0, 64)
	threshold.Value = 10
	live := widget.NewCheck("● Тест в реальном времени", func(v bool) { r.liveTest = v })
	btnSave := widget.NewButton("Сохранить область →", func() {
		if r.OnSave != nil {
			r.OnSave(r.selection, int(math.Round(threshold.Value)), "")
		}
	})

	r.root = container.NewBorder(
		container.NewVBox(widget.NewLabel("Редактор области")),
		container.NewHBox(widget.NewButton("← Отмена", func() {}), btnSave),
		layer,
		container.NewVBox(
			widget.NewLabel("📐 Координаты"), x, y, w, h,
			widget.NewLabel("⚙️ Порог"), threshold,
			live,
		),
	)
}

func (r *RegionEditor) Widget() fyne.CanvasObject         { return r.root }
func (r *RegionEditor) SetSelection(rect image.Rectangle) { r.selection = rect }
func (r *RegionEditor) widgetToImage(p fyne.Position) image.Point {
	return image.Point{X: int((p.X - r.pan.X) / r.zoom), Y: int((p.Y - r.pan.Y) / r.zoom)}
}
