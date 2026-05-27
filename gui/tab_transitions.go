package gui

import (
	"fmt"
	"sort"
	"vpmc/config"
	"vpmc/gui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func buildTransitionsTab(cfg *config.Config) *container.TabItem {
	stateIDs := make([]string, 0, len(cfg.States))
	for id := range cfg.States {
		stateIDs = append(stateIDs, id)
	}
	sort.Strings(stateIDs)

	rows := container.NewVBox(widget.NewLabel("Выберите состояние"))
	scroll := container.NewVScroll(rows)
	sel := widget.NewSelect(stateIDs, func(stateID string) {
		objs := []fyne.CanvasObject{widget.NewLabel("Переходы: " + stateID)}
		for _, tr := range cfg.States[stateID].Transitions {
			trLocal := tr
			desc := fmt.Sprintf("%s → %s (thr=%d, refs=%d, %s)", trLocal.ID, trLocal.ToState, trLocal.Threshold, len(trLocal.References), rectToString(trLocal.Rect))
			objs = append(objs, widget.NewButton(desc, func() {
				editor := widgets.NewRegionEditor(nil)
				editor.SetSelection(trLocal.Rect)
				w := fyne.CurrentApp().Driver().AllWindows()[0]
				d := dialog.NewCustom("Редактор области", "Закрыть", editor.Widget(), w)
				d.Resize(fyne.NewSize(1024, 700))
				d.Show()
			}))
		}
		rows.Objects = objs
		rows.Refresh()
	})
	sel.PlaceHolder = "Состояние"

	content := container.NewBorder(container.NewVBox(widget.NewLabel("Переходы"), sel), nil, nil, nil, scroll)
	return container.NewTabItem("Переходы", content)
}
