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
	mode := widget.NewRadioGroup([]string{"auto", "winkey", "pear"}, func(s string) {
		cfg.ControllerMode = s
		mgr.SetMode(controller.Mode(s))
		_ = config.Save("config.json", cfg)
	})
	mode.Selected = cfg.ControllerMode
	port := widget.NewEntry()
	port.SetText(strconv.Itoa(cfg.PearPort))
	authID := widget.NewEntry()
	authID.SetText(cfg.PearAuthID)
	token := widget.NewPasswordEntry()
	token.SetText(cfg.PearToken)
	status := widget.NewLabel("Статус контроллеров")

	reload := func() {
		st := mgr.StatusAll()
		status.SetText(fmt.Sprintf("WinKey: %v\nPear: %v", st[controller.ModeWinKey].Available, st[controller.ModePear].Available))
	}
	checkBtn := widget.NewButton("Проверить", func() {
		if p, err := strconv.Atoi(port.Text); err == nil {
			cfg.PearPort = p
			mgr.SetPearPort(p)
			_ = config.Save("config.json", cfg)
		}
		cfg.PearAuthID = authID.Text
		cfg.PearToken = token.Text
		mgr.SetPearConfig(cfg.PearPort, cfg.PearToken, cfg.PearAuthID)
		_ = config.Save("config.json", cfg)
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
		widget.NewButton("▶ Play", func() { status.SetText(fmt.Sprintf("Play: %v", mgr.Play())); reload() }),
		widget.NewButton("⏸ Pause", func() { status.SetText(fmt.Sprintf("Pause: %v", mgr.Pause())); reload() }),
		widget.NewButton("⟳ Toggle", func() { status.SetText(fmt.Sprintf("Toggle: %v", mgr.Toggle())); reload() }),
	)

	content := container.NewVBox(
		widget.NewLabel("🎵 Режим управления музыкой"),
		mode,
		container.NewHBox(widget.NewLabel("Порт Pear:"), port, checkBtn),
		status,
		testBtns,
	)
	return container.NewTabItem("Плеер", content)
}
