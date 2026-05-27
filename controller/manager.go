package controller

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"vpmc/config"
)

type Mode string

const (
	ModeWinKey Mode = "winkey"
	ModePear   Mode = "pear"
	ModeAuto   Mode = "auto"
)

type Status struct {
	Name      string
	Available bool
}

type Manager struct {
	controllers     map[Mode]MusicController
	mu              sync.RWMutex
	active          atomic.Value
	lastKnownAction atomic.Value
	logger          *slog.Logger
}

func NewManager(cfg *config.Config) *Manager {
	m := &Manager{controllers: map[Mode]MusicController{}, logger: slog.Default()}
	m.controllers[ModeWinKey] = NewWinKeyController()
	m.controllers[ModePear] = NewPearController(cfg.PearPort)
	m.active.Store(Mode(cfg.ControllerMode))
	m.lastKnownAction.Store(config.ActionPlay)
	return m
}

func (m *Manager) activeCtl() MusicController {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a := m.active.Load().(Mode)
	if a == ModeAuto {
		if p, ok := m.controllers[ModePear]; ok && p.IsAvailable() {
			return p
		}
		return m.controllers[ModeWinKey]
	}
	if c, ok := m.controllers[a]; ok {
		return c
	}
	return m.controllers[ModeWinKey]
}

func (m *Manager) SetMode(mode Mode) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.controllers[mode]; !ok && mode != ModeAuto {
		mode = ModeAuto
	}
	m.active.Store(mode)
}

func (m *Manager) SetPearPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.controllers[ModePear] = NewPearController(port)
}
func (m *Manager) Play() error {
	m.lastKnownAction.Store(config.ActionPlay)
	return m.activeCtl().Play()
}
func (m *Manager) Pause() error {
	m.lastKnownAction.Store(config.ActionPause)
	return m.activeCtl().Pause()
}
func (m *Manager) Toggle() error { return m.activeCtl().Toggle() }
func (m *Manager) ExecuteAction(a config.Action) {
	switch a {
	case config.ActionPlay:
		_ = m.Play()
	case config.ActionPause:
		_ = m.Pause()
	}
}
func (m *Manager) StatusAll() map[Mode]Status {
	out := map[Mode]Status{}
	for mode, c := range m.controllers {
		out[mode] = Status{Name: c.Name(), Available: c.IsAvailable()}
	}
	return out
}
