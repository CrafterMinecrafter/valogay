package gui

import (
	"strconv"
	"vpmc/config"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildSettingsTab(cfg *config.Config) *container.TabItem {
	base := widget.NewEntry()
	base.SetText(strconv.Itoa(cfg.Monitor.BaseIntervalMs))
	hys := widget.NewEntry()
	hys.SetText(strconv.Itoa(cfg.Monitor.HysteresisCount))
	watch := widget.NewEntry()
	watch.SetText(strconv.Itoa(cfg.Monitor.WatchdogTimeoutSec))
	mode := widget.NewEntry()
	mode.SetText(strconv.Itoa(cfg.Monitor.ModeDetectThreshold))
	save := widget.NewButton("Применить", func() {
		if v, e := strconv.Atoi(base.Text); e == nil {
			cfg.Monitor.BaseIntervalMs = v
		}
		if v, e := strconv.Atoi(hys.Text); e == nil {
			cfg.Monitor.HysteresisCount = v
		}
		if v, e := strconv.Atoi(watch.Text); e == nil {
			cfg.Monitor.WatchdogTimeoutSec = v
		}
		if v, e := strconv.Atoi(mode.Text); e == nil {
			cfg.Monitor.ModeDetectThreshold = v
		}
		_ = config.Save("config.json", cfg)
	})
	return container.NewTabItem("Настройки", container.NewVBox(widget.NewLabel("Monitor"), base, hys, watch, mode, save))
}
