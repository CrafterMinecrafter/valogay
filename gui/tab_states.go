package gui

import (
	"fmt"
	"sort"
	"vpmc/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildStatesTab(cfg *config.Config) *container.TabItem {
	ids := make([]string, 0, len(cfg.States))
	for id := range cfg.States {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	list := widget.NewList(
		func() int { return len(ids) },
		func() fyne.CanvasObject { return widget.NewLabel("state") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			id := ids[i]
			s := cfg.States[id]
			o.(*widget.Label).SetText(fmt.Sprintf("%s (%s) • action=%s • interval=%dms", id, s.Name, s.Action, s.IntervalMs))
		},
	)

	head := widget.NewRichTextFromMarkdown("### Состояния FSM\nСписок состояний и базовых действий.")
	return container.NewTabItem("Состояния", container.NewBorder(head, nil, nil, nil, list))
}
