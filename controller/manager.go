package controller

import (
	"log/slog"
	"sync/atomic"
	"vpmc/config"
)

type Mode string
const (ModeWinKey Mode="winkey"; ModePear Mode="pear"; ModeAuto Mode="auto")

type Manager struct { controllers map[Mode]MusicController; active atomic.Value; lastKnownAction atomic.Value; logger *slog.Logger }
func NewManager(cfg *config.Config)*Manager{m:=&Manager{controllers:map[Mode]MusicController{},logger:slog.Default()}; m.controllers[ModeWinKey]=NewWinKeyController(); m.active.Store(Mode(cfg.ControllerMode)); m.lastKnownAction.Store(config.ActionPlay); return m}
func (m *Manager) activeCtl() MusicController { a:=m.active.Load().(Mode); if a==ModeAuto {if p,ok:=m.controllers[ModePear];ok&&p.IsAvailable(){return p}; return m.controllers[ModeWinKey]}; if c,ok:=m.controllers[a];ok{return c}; return m.controllers[ModeWinKey] }
func (m *Manager) Play() error {m.lastKnownAction.Store(config.ActionPlay); return m.activeCtl().Play()}
func (m *Manager) Pause() error {m.lastKnownAction.Store(config.ActionPause); return m.activeCtl().Pause()}
func (m *Manager) ExecuteAction(a config.Action){switch a{case config.ActionPlay:_=m.Play();case config.ActionPause:_=m.Pause()}}
