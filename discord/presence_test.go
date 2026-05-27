package discord

import (
	"testing"
	"vpmc/config"
)

func TestButtons(t *testing.T) {
	p := NewPresenceManager(&config.DiscordConfig{RiotID: "Nick#RU1", CustomBtn: &config.DiscordButton{Label: "Site", URL: "https://example.com"}})
	b := p.buildButtons()
	if len(b) != 2 {
		t.Fatalf("want 2 got %d", len(b))
	}
}
