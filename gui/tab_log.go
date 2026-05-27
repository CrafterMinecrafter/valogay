package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildLogTab() *container.TabItem {
	logBox := widget.NewMultiLineEntry()
	logBox.SetPlaceHolder("Лог мониторинга будет отображаться здесь")
	logBox.Disable()
	return container.NewTabItem("Лог", container.NewMax(logBox))
}
