package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"vpmc/config"
	"vpmc/controller"
	"vpmc/discord"
	"vpmc/fsm"
	"vpmc/monitor"
	"vpmc/recorder"
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
	rec      *recorder.Recorder
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

func (a *App) GetRegions() []config.Region { return a.cfg.Regions }

func (a *App) SaveRegions(regions []config.Region) error {
	a.cfg.Regions = regions
	return config.Save("config.json", a.cfg)
}

func (a *App) DeleteRegion(id string) error {
	regions := make([]config.Region, 0, len(a.cfg.Regions))
	for _, r := range a.cfg.Regions {
		if r.ID != id {
			regions = append(regions, r)
		}
	}
	a.cfg.Regions = regions
	return config.Save("config.json", a.cfg)
}

func (a *App) GetStates() map[string]config.StateConfig { return a.cfg.States }

func (a *App) RecorderStart() error {
	if a.rec != nil {
		a.cfg.Recorder.Enabled = true
		a.rec.Start()
	}
	return nil
}

func (a *App) RecorderStop() error {
	if a.rec != nil {
		a.cfg.Recorder.Enabled = false
		a.rec.Stop()
	}
	return nil
}

func (a *App) RecorderStatus() map[string]interface{} {
	count := 0
	if a.rec != nil {
		count = len(a.rec.Frames())
	}
	return map[string]interface{}{
		"enabled":      a.cfg.Recorder.Enabled,
		"interval_sec": a.cfg.Recorder.IntervalSec,
		"output_dir":   a.cfg.Recorder.OutputDir,
		"frames_count": count,
	}
}

func (a *App) RecorderLastFrameImage() string {
	if a.rec == nil {
		return ""
	}
	f := a.rec.LastFrame()
	if f == nil {
		return ""
	}
	img := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			v := f.Pixels[y*f.Width+x]
			img.Set(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (a *App) RecorderFramesMeta() []map[string]interface{} {
	if a.rec == nil {
		return nil
	}
	frames := a.rec.Frames()
	out := make([]map[string]interface{}, len(frames))
	for i, f := range frames {
		out[i] = map[string]interface{}{
			"width":     f.Width,
			"height":    f.Height,
			"timestamp": f.Timestamp,
		}
	}
	return out
}


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

	rec := recorder.New(cfg.Recorder.IntervalSec, cfg.Recorder.OutputDir)
	if cfg.Recorder.Enabled {
		rec.Start()
	}

	app := &App{cfg: cfg, machine: machine, mgr: mgr, presence: presence, mon: monitor.New(machine, mgr, vision.NewComparer(), capturer, cfg), capturer: capturer, rec: rec}

	go app.mon.Start()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		_ = mgr.Play()
		presence.Logout()
		app.mon.Stop()
		rec.Stop()
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
