package monitor

import (
	"context"
	"image"
	"log/slog"
	"sync"
	"time"
	"vpmc/config"
	"vpmc/controller"
	"vpmc/fsm"
	"vpmc/vision"
)

type Monitor struct {
	fsm        *fsm.FSM
	controller *controller.Manager
	vision     *vision.Comparer
	capturer   vision.Capturer
	cfg        *config.Config
	hysteresis map[string]*Hysteresis
	stopCh     chan struct{}
	wg         sync.WaitGroup
	watchdog   *Watchdog
}

func New(machine *fsm.FSM, mgr *controller.Manager, cmp *vision.Comparer, capturer vision.Capturer, cfg *config.Config) *Monitor {
	m := &Monitor{fsm: machine, controller: mgr, vision: cmp, capturer: capturer, cfg: cfg, stopCh: make(chan struct{}), hysteresis: map[string]*Hysteresis{}, watchdog: NewWatchdog(cfg.Monitor.WatchdogTimeoutSec)}
	for _, st := range cfg.States {
		for _, tr := range st.Transitions {
			m.hysteresis[tr.ID] = NewHysteresis(cfg.Monitor.HysteresisCount)
		}
	}
	return m
}

func intersectsAny(target image.Rectangle, dirty []image.Rectangle) bool {
	if len(dirty) == 0 {
		return true
	}
	for _, d := range dirty {
		if target.Overlaps(d) {
			return true
		}
	}
	return false
}

func (m *Monitor) Start() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				slog.Error("monitor panic", "recover", r)
			}
		}()
		for {
			select {
			case <-m.stopCh:
				return
			default:
			}
			st := m.fsm.State()
			if st == nil {
				time.Sleep(time.Second)
				continue
			}
			if m.fsm.Current() == "MODE_DETECT" {
				_ = m.handleModeDetect(context.Background())
			}
			dirty := m.capturer.DirtyRects()
			for _, tr := range st.Transitions {
				if !intersectsAny(tr.Rect, dirty) {
					continue
				}
				img, err := m.capturer.CaptureRect(tr.Rect)
				if err != nil {
					continue
				}
				ok, err := m.vision.MatchesAny(img, tr.References, tr.Threshold)
				if err != nil || !ok {
					m.hysteresis[tr.ID].Reset()
					continue
				}
				if m.hysteresis[tr.ID].Hit() {
					a := tr.Action
					if a == "" {
						a = m.fsm.ResolveAction(tr.ToState)
					}
					m.fsm.Transition(tr.ToState, a)
					m.watchdog.Touch()
					break
				}
			}
			if m.watchdog.Expired() && m.fsm.Current() != "LAUNCHER" {
				m.fsm.Transition("LAUNCHER", config.ActionPlay)
				_ = m.controller.Play()
				m.watchdog.Touch()
			}
			interval := m.cfg.Monitor.BaseIntervalMs
			if st.IntervalMs > 0 {
				interval = st.IntervalMs
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
}

func (m *Monitor) handleModeDetect(_ context.Context) error {
	for _, mode := range m.cfg.GameModes {
		img, err := m.capturer.CaptureRect(mode.DetectRect)
		if err != nil {
			continue
		}
		ok, err := m.vision.MatchesAny(img, []string{mode.DetectImage}, m.cfg.Monitor.ModeDetectThreshold)
		if err == nil && ok {
			m.fsm.SetMode(mode.ID)
			if mode.ID == "deathmatch" {
				m.fsm.Transition("IN_ALIVE", m.fsm.ResolveAction("IN_ALIVE"))
			} else {
				m.fsm.Transition("LOADING", m.fsm.ResolveAction("LOADING"))
			}
			return nil
		}
	}
	m.fsm.SetMode("unrated")
	m.fsm.Transition("LOADING", m.fsm.ResolveAction("LOADING"))
	return nil
}

func (m *Monitor) Stop() { close(m.stopCh); m.wg.Wait() }
