package discord

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"
	"vpmc/config"
)

type Button struct{ Label, URL string }
type Activity struct {
	Details, State string
	Buttons        []Button
}

type PresenceManager struct {
	cfg          *config.DiscordConfig
	connected    atomic.Bool
	sessionStart time.Time
	roundStart   time.Time
	lastUpdate   time.Time
	mu           sync.Mutex
	lastActivity Activity
}

func NewPresenceManager(cfg *config.DiscordConfig) *PresenceManager {
	return &PresenceManager{cfg: cfg, sessionStart: time.Now()}
}
func (p *PresenceManager) Start()  {}
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
	return b
}
func (p *PresenceManager) buildActivity(state, mode string, round int) Activity {
	return Activity{Details: state, State: mode, Buttons: p.buildButtons()}
}
