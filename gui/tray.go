package gui

import (
	"fmt"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
	"vpmc/gui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func buildTray(a fyne.App, w fyne.Window, machine *fsm.FSM, manager *controller.Manager, presence *discord.PresenceManager, _ *widgets.StatusBar) (*fyne.MenuItem, *fyne.MenuItem) {
	stateItem := fyne.NewMenuItem("Состояние: "+machine.Current(), nil)
	mode := machine.CurrentMode()
	if mode == "" {
		mode = "—"
	}
	modeItem := fyne.NewMenuItem("Режим: "+mode, nil)

	showMainWindow := func() { w.Show(); w.RequestFocus() }
	onExit := func() {
		_ = manager.Play()
		presence.Logout()
		w.Close()
		a.Quit()
	}

	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayMenu(fyne.NewMenu("VPMC",
			fyne.NewMenuItem("Открыть VPMC", showMainWindow),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("▶ Play", func() { _ = manager.Play() }),
			fyne.NewMenuItem("⏸ Pause", func() { _ = manager.Pause() }),
			fyne.NewMenuItemSeparator(),
			stateItem,
			modeItem,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem(fmt.Sprintf("Выход"), onExit),
		))
	}
	return stateItem, modeItem
}
