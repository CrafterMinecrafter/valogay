package gui

import (
	"fmt"
	"image"
	"sync"
	"vpmc/config"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
	"vpmc/gui/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type AppUI struct {
	App       fyne.App
	Win       fyne.Window
	StateBar  *widgets.StatusBar
	stateMu   sync.Mutex
	stateItem *fyne.MenuItem
	modeItem  *fyne.MenuItem
}

func Run(cfg *config.Config, machine *fsm.FSM, mgr *controller.Manager, presence *discord.PresenceManager) {
	a := app.NewWithID("vpmc")
	w := a.NewWindow("VPMC — Valorant Phase Music Controller")
	w.Resize(fyne.NewSize(1180, 760))

	ui := &AppUI{App: a, Win: w, StateBar: widgets.NewStatusBar()}
	tabs := container.NewAppTabs(
		buildStatesTab(cfg),
		buildTransitionsTab(cfg),
		buildControllerTab(cfg, mgr),
		buildDiscordTab(cfg, presence),
		buildLogTab(),
		buildSettingsTab(cfg),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	ui.bindFSM(machine)
	ui.stateItem, ui.modeItem = buildTray(a, w, machine, mgr, presence, ui.StateBar)
	w.SetContent(container.NewBorder(nil, ui.StateBar.Widget(), nil, nil, tabs))
	w.SetCloseIntercept(func() { w.Hide() })
	w.ShowAndRun()
}

func (u *AppUI) bindFSM(machine *fsm.FSM) {
	prev := machine.OnTransition
	machine.OnTransition = func(from, to string, action config.Action) {
		if prev != nil {
			prev(from, to, action)
		}
		u.stateMu.Lock()
		defer u.stateMu.Unlock()
		mode := machine.CurrentMode()
		u.StateBar.Update(to, mode, string(action), machine.RoundNum())
		if u.stateItem != nil {
			u.stateItem.Label = fmt.Sprintf("Состояние: %s", to)
		}
		if u.modeItem != nil {
			if mode == "" {
				mode = "—"
			}
			u.modeItem.Label = fmt.Sprintf("Режим: %s", mode)
		}
	}
}

func rectToString(r image.Rectangle) string {
	return fmt.Sprintf("X:%d Y:%d W:%d H:%d", r.Min.X, r.Min.Y, r.Dx(), r.Dy())
}
