//go:build windows

package recorder

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	gdi32  = windows.NewLazySystemDLL("gdi32.dll")

	procGetDC            = user32.NewProc("GetDC")
	procReleaseDC        = user32.NewProc("ReleaseDC")
	procCreateDC         = gdi32.NewProc("CreateDCW")
	procCreateCompatibleDC   = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject     = gdi32.NewProc("SelectObject")
	procBitBlt           = gdi32.NewProc("BitBlt")
	procGetDeviceCaps    = gdi32.NewProc("GetDeviceCaps")
	procDeleteDC         = gdi32.NewProc("DeleteDC")
	procDeleteObject     = gdi32.NewProc("DeleteObject")
	procGetDIBits        = gdi32.NewProc("GetDIBits")
)

const (
	SRCCOPY = 0x00CC0020
	CAPTUREBLT = 0x40000000
	DESKTOPHORZRES = 118
	DESKTOPVERTRES = 117
)

type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type BITMAPINFO struct {
	Header BITMAPINFOHEADER
	Colors [1]uint32
}

func captureScreen() (*image.RGBA, error) {
	hdcScreen, _, _ := procGetDC.Call(0)
	if hdcScreen == 0 {
		return nil, fmt.Errorf("GetDC failed")
	}
	defer procReleaseDC.Call(0, hdcScreen)

	width := getDeviceCaps(hdcScreen, DESKTOPHORZRES)
	height := getDeviceCaps(hdcScreen, DESKTOPVERTRES)

	hdcMem, _, _ := procCreateCompatibleDC.Call(hdcScreen)
	if hdcMem == 0 {
		return nil, fmt.Errorf("CreateCompatibleDC failed")
	}
	defer procDeleteDC.Call(hdcMem)

	hBitmap, _, _ := procCreateCompatibleBitmap.Call(hdcScreen, uintptr(width), uintptr(height))
	if hBitmap == 0 {
		return nil, fmt.Errorf("CreateCompatibleBitmap failed")
	}
	defer procDeleteObject.Call(hBitmap)

	oldObj, _, _ := procSelectObject.Call(hdcMem, hBitmap)
	if oldObj == 0 {
		return nil, fmt.Errorf("SelectObject failed")
	}
	defer procSelectObject.Call(hdcMem, oldObj)

	ret, _, _ := procBitBlt.Call(hdcMem, 0, 0, uintptr(width), uintptr(height), hdcScreen, 0, 0, SRCCOPY|CAPTUREBLT)
	if ret == 0 {
		return nil, fmt.Errorf("BitBlt failed")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var bi BITMAPINFOHEADER
	bi.Size = uint32(unsafe.Sizeof(bi))
	bi.Width = int32(width)
	bi.Height = -int32(height)
	bi.Planes = 1
	bi.BitCount = 32
	bi.Compression = 0

	bmi := &BITMAPINFO{Header: bi}

	ret, _, _ = procGetDIBits.Call(hdcMem, hBitmap, 0, uintptr(height),
		uintptr(unsafe.Pointer(&img.Pix[0])),
		uintptr(unsafe.Pointer(bmi)),
		0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetDIBits failed")
	}

	return img, nil
}

func getDeviceCaps(hdc uintptr, index int) int {
	ret, _, _ := procGetDeviceCaps.Call(hdc, uintptr(index))
	return int(ret)
}

func saveFrame(img *image.RGBA, dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	name := filepath.Join(dir, fmt.Sprintf("valo_%s.jpg", time.Now().Format("20060102_150405")))
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()
	jpeg.Encode(f, img, &jpeg.Options{Quality: 80})
}
