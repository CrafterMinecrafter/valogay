package config

import (
	"encoding/json"
	"os"
)

func Default() *Config {
	return Migrate(&Config{
		Version: 1, ControllerMode: "auto",
		PearPort: 9863, PearAuthID: "", PearToken: "",
		States:    map[string]StateConfig{},
		GameModes: map[string]GameMode{},
		Monitor: MonitorConfig{
			BaseIntervalMs:      500,
			HysteresisCount:     3,
			WatchdogTimeoutSec:  30,
			ModeDetectThreshold: 10,
		},
		Recorder: RecorderConfig{
			Enabled:     false,
			IntervalSec: 5,
			OutputDir:   "recordings",
		},
		Discord: DiscordConfig{},
	})
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return Migrate(&c), nil
}

func Save(path string, cfg *Config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
