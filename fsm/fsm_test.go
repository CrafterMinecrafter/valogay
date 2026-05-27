package fsm

import (
	"testing"
	"vpmc/config"
)

func TestResolveActionOverride(t *testing.T) {
	cfg := &config.Config{States: map[string]config.StateConfig{"IN_ALIVE": {Action: config.ActionPause}}, GameModes: map[string]config.GameMode{"deathmatch": {ID: "deathmatch", ActionOverrides: map[string]config.Action{"IN_ALIVE": config.ActionPlay}}}}
	m := New(cfg)
	m.SetMode("deathmatch")
	if got := m.ResolveAction("IN_ALIVE"); got != config.ActionPlay {
		t.Fatalf("want play got %s", got)
	}
}
