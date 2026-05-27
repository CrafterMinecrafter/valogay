package discord

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
	"vpmc/config"
)

type Button struct{ Label, URL string }
type Activity struct {
	Details, State, LargeImage string
	Buttons                    []Button
	Start                      time.Time
}

type PresenceManager struct {
	cfg          *config.DiscordConfig
	connected    atomic.Bool
	sessionStart time.Time
	roundStart   time.Time
	lastUpdate   time.Time
	mu           sync.Mutex
	stopCh       chan struct{}
	lastActivity Activity
}

func NewPresenceManager(cfg *config.DiscordConfig) *PresenceManager {
	return &PresenceManager{cfg: cfg, sessionStart: time.Now(), stopCh: make(chan struct{})}
}

func (p *PresenceManager) Start(ctx context.Context) {
	go func() { defer func() { recover() }(); <-ctx.Done(); p.Logout() }()
}
func (p *PresenceManager) Logout() { p.connected.Store(false) }

func (p *PresenceManager) Update(state, mode string, round int) error {
	if !p.cfg.Enabled {
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if time.Since(p.lastUpdate) < 15*time.Second {
		return nil
	}
	if state == "LAUNCHER" {
		p.connected.Store(false)
		p.lastUpdate = time.Now()
		return nil
	}
	p.connected.Store(true)
	p.lastActivity = p.buildActivity(state, mode, round)
	p.lastUpdate = time.Now()
	return nil
}
func (p *PresenceManager) buildButtons() []Button {
	var b []Button
	if p.cfg.RiotID != "" {
		b = append(b, Button{"Профиль на Tracker.gg", "https://tracker.gg/valorant/profile/riot/" + url.PathEscape(p.cfg.RiotID)})
	}
	if p.cfg.CustomBtn != nil && p.cfg.CustomBtn.Label != "" {
		b = append(b, Button{p.cfg.CustomBtn.Label, p.cfg.CustomBtn.URL})
	}
	if len(b) > 2 {
		return b[:2]
	}
	if len(b) == 0 {
		return nil
	}
	return b
}
func (p *PresenceManager) buildActivity(state, mode string, round int) Activity {
	act := Activity{State: mode, Buttons: p.buildButtons(), Start: p.sessionStart}
	switch state {
	case "MAIN_MENU":
		act.LargeImage = "val_mainmenu"
		act.Details = "В главном меню"
	case "AGENT_SELECT", "MODE_DETECT":
		act.LargeImage = "val_agentselect"
		act.Details = "Выбор агента"
	case "LOADING":
		act.LargeImage = "val_agentselect"
		act.Details = "Загрузка матча"
		act.Start = time.Now()
	case "IN_BUY":
		act.LargeImage = "val_buy"
		act.Details = fmt.Sprintf("Фаза закупки • Раунд %d", round)
		if p.roundStart.IsZero() {
			p.roundStart = time.Now()
		}
		act.Start = p.roundStart
	case "IN_ALIVE":
		if mode == "deathmatch" {
			act.LargeImage = "val_dm"
			act.Details = "Deathmatch"
		} else {
			act.LargeImage = "val_alive"
			act.Details = fmt.Sprintf("Раунд %d", round)
		}
		if !p.roundStart.IsZero() {
			act.Start = p.roundStart
		}
	case "IN_DEAD":
		if mode == "deathmatch" {
			act.LargeImage = "val_dm"
			act.Details = "Deathmatch"
		} else {
			act.LargeImage = "val_dead"
			act.Details = fmt.Sprintf("Раунд %d", round)
		}
		if !p.roundStart.IsZero() {
			act.Start = p.roundStart
		}
	case "MATCH_END":
		act.LargeImage = "val_matchend"
		act.Details = "Матч завершён"
	}
	return act
}
