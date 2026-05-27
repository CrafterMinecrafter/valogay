package main

import (
	"embed"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vpmc/config"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
	"vpmc/monitor"
	"vpmc/vision"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

type App struct {
	cfg      *config.Config
	machine  *fsm.FSM
	mgr      *controller.Manager
	presence *discord.PresenceManager
	mon      *monitor.Monitor
	capturer vision.Capturer
}

type Status struct {
	WinKey bool `json:"WinKey"`
	Pear   bool `json:"Pear"`
}

type SongInfo struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	IsPaused bool   `json:"isPaused"`
}

func (a *App) Play() error         { return a.mgr.Play() }
func (a *App) Pause() error        { return a.mgr.Pause() }
func (a *App) Toggle() error       { return a.mgr.Toggle() }
func (a *App) GetStatus() Status    { return Status{WinKey: a.mgr.StatusAll()[controller.ModeWinKey].Available, Pear: a.mgr.StatusAll()[controller.ModePear].Available} }

func (a *App) GetSongInfo() *SongInfo {
	if pc, ok := a.mgr.Controller().(interface{ SongInfo() (map[string]interface{}, error) }); ok {
		if info, err := pc.SongInfo(); err == nil && info != nil {
			return &SongInfo{Title: getVal[string](info, "title"), Artist: getVal[string](info, "artist"), IsPaused: getVal[bool](info, "isPaused")}
		}
	}
	return nil
}

func (a *App) GetStateInfo() map[string]interface{} {
	return map[string]interface{}{"state": a.machine.Current(), "mode": a.machine.CurrentMode(), "action": string(a.mgr.GetLastAction()), "round": a.machine.RoundNum()}
}

func getVal[T any](m map[string]interface{}, key string) T {
	if v, ok := m[key].(T); ok {
		return v
	}
	var zero T
	return zero
}

func (a *App) GetConfig() *config.Config                { return a.cfg }
func (a *App) SetMode(mode string)                      { a.cfg.ControllerMode = mode; a.mgr.SetMode(controller.Mode(mode)) }
func (a *App) SetPearPort(port int)                     { a.cfg.PearPort = port; a.mgr.SetPearConfig(port) }
func (a *App) Start()                                  { a.mon.Start() }
func (a *App) Stop()                                   { a.mon.Stop() }
func (a *App) GetDiscordConfig() *config.DiscordConfig { return &a.cfg.Discord }

func (a *App) SaveDiscordConfig(enabled bool, appID, riotID, btnLabel, cmdURL string) error {
	a.cfg.Discord.Enabled = enabled
	a.cfg.Discord.AppID = appID
	a.cfg.Discord.RiotID = riotID
	if btnLabel != "" {
		a.cfg.Discord.CustomBtn = &config.DiscordButton{Label: btnLabel, URL: cmdURL}
	}
	return config.Save("config.json", a.cfg)
}

func (a *App) GetMonitorConfig() *config.MonitorConfig { return &a.cfg.Monitor }

func (a *App) SaveMonitorConfig(intervalMs, hysteresis, watchdogSec, modeThreshold int) error {
	a.cfg.Monitor.BaseIntervalMs = intervalMs
	a.cfg.Monitor.HysteresisCount = hysteresis
	a.cfg.Monitor.WatchdogTimeoutSec = watchdogSec
	a.cfg.Monitor.ModeDetectThreshold = modeThreshold
	return config.Save("config.json", a.cfg)
}

func (a *App) GetStates() map[string]config.StateConfig { return a.cfg.States }

func main() {
	cfg, _ := config.Load("config.json")
	if cfg == nil {
		cfg = config.Default()
	}

	capturer, err := vision.NewDXGICapturer(cfg.Monitor.DisplayIndex)
	if err != nil {
		slog.Error("capturer init", "err", err)
		return
	}
	defer capturer.Close()

	machine := fsm.New(cfg)
	machine.Reset()
	mgr := controller.NewManager(cfg)
	presence := discord.NewPresenceManager(&cfg.Discord)

	machine.OnTransition = func(_ string, to string, action config.Action) {
		mgr.ExecuteAction(action)
		_ = presence.Update(to, machine.CurrentMode(), machine.RoundNum())
	}

	app := &App{cfg: cfg, machine: machine, mgr: mgr, presence: presence, mon: monitor.New(machine, mgr, vision.NewComparer(), capturer, cfg), capturer: capturer}

	go app.mon.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = mgr.Play()
		presence.Logout()
		app.mon.Stop()
		_ = capturer.Close()
		os.Exit(0)
	}()

	wails.Run(&options.App{
		Title:           "VPMC — Valorant Phase Music Controller",
		Width:           900,
		Height:          650,
		AssetServer:     &assetserver.Options{Assets: assets},
		BackgroundColour: &options.RGBA{R: 15, G: 15, B: 15, A: 255},
		Bind:            []interface{}{app},
	})
}
