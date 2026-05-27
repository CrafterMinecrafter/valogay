package gui

import (
	"fmt"
	"strconv"
	"time"
	"vpmc/config"
	"vpmc/controller"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildControllerTab(cfg *config.Config, mgr *controller.Manager) *container.TabItem {
	mode := widget.NewRadioGroup([]string{"auto", "winkey", "pear"}, func(s string) { cfg.ControllerMode = s })
	mode.Selected = cfg.ControllerMode
	port := widget.NewEntry()
	port.SetText(strconv.Itoa(cfg.PearPort))
	status := widget.NewLabel("Статус контроллеров")

	reload := func() {
		st := mgr.StatusAll()
		status.SetText(fmt.Sprintf("WinKey: %v\nPear: %v", st[controller.ModeWinKey].Available, st[controller.ModePear].Available))
	}
	widget.NewButton("Проверить", func() {
		if p, err := strconv.Atoi(port.Text); err == nil {
			cfg.PearPort = p
		}
		reload()
	})
	reload()

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for range t.C {
			reload()
		}
	}()

	testBtns := container.NewHBox(
		widget.NewButton("▶ Play", func() { _ = mgr.Play() }),
		widget.NewButton("⏸ Pause", func() { _ = mgr.Pause() }),
		widget.NewButton("⟳ Toggle", func() { _ = mgr.Toggle() }),
	)

	content := container.NewVBox(
		widget.NewLabel("🎵 Режим управления музыкой"),
		mode,
		container.NewHBox(widget.NewLabel("Порт Pear:"), port, widget.NewButton("Проверить", func() { reload() })),
		status,
		testBtns,
	)
	return container.NewTabItem("Плеер", content)
}
