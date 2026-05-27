package monitor

import (
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
			interval := m.cfg.Monitor.BaseIntervalMs
			if st.IntervalMs > 0 {
				interval = st.IntervalMs
			}
			if m.watchdog.Expired() && m.fsm.Current() != "LAUNCHER" {
				m.fsm.Transition("LAUNCHER", config.ActionPlay)
				_ = m.controller.Play()
				m.watchdog.Touch()
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
}

func (m *Monitor) Stop() { close(m.stopCh); m.wg.Wait() }
