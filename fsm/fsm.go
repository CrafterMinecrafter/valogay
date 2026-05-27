package fsm

import (
	"sync"
	"vpmc/config"
)

type FSM struct {
	states map[string]*config.StateConfig
	gameModes map[string]*config.GameMode
	current string
	currentMode *config.GameMode
	roundNum int
	mu sync.RWMutex
	OnTransition func(from,to string,action config.Action)
	OnModeSet func(mode *config.GameMode)
}

func New(cfg *config.Config) *FSM {
	s:=map[string]*config.StateConfig{}
	for k,v:= range cfg.States { vv:=v; s[k]=&vv }
	gm:=map[string]*config.GameMode{}
	for k,v:= range cfg.GameModes { vv:=v; gm[k]=&vv }
	return &FSM{states:s,gameModes:gm,current:"LAUNCHER"}
}
func (f *FSM) Current() string { f.mu.RLock(); defer f.mu.RUnlock(); return f.current }
func (f *FSM) CurrentMode() string { f.mu.RLock(); defer f.mu.RUnlock(); if f.currentMode==nil{return ""}; return f.currentMode.ID }
func (f *FSM) RoundNum() int { f.mu.RLock(); defer f.mu.RUnlock(); return f.roundNum }
func (f *FSM) State() *config.StateConfig { f.mu.RLock(); defer f.mu.RUnlock(); return f.states[f.current] }
func (f *FSM) Transition(to string, action config.Action) {
	f.mu.Lock(); from:=f.current; f.current=to; if to=="IN_BUY" {f.roundNum++}; f.mu.Unlock()
	if f.OnTransition!=nil { f.OnTransition(from,to,action)}
}
func (f *FSM) SetMode(id string) { f.mu.Lock(); m:=f.gameModes[id]; f.currentMode=m; f.mu.Unlock(); if m!=nil && f.OnModeSet!=nil {f.OnModeSet(m)} }
func (f *FSM) ResolveAction(stateID string) config.Action { f.mu.RLock(); defer f.mu.RUnlock(); if f.currentMode!=nil { if o,ok:=f.currentMode.ActionOverrides[stateID]; ok {return o} }; if s:=f.states[stateID]; s!=nil {return s.Action}; return config.ActionNone }
func (f *FSM) Reset() { f.mu.Lock(); defer f.mu.Unlock(); f.current="MAIN_MENU"; f.currentMode=nil; f.roundNum=0 }
