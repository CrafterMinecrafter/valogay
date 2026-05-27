package config

func Migrate(cfg *Config) *Config {
	if cfg == nil {
		return Default()
	}
	if cfg.Version <= 0 {
		cfg.Version = 1
	}
	if cfg.States == nil {
		cfg.States = map[string]StateConfig{}
	}
	if cfg.GameModes == nil {
		cfg.GameModes = map[string]GameMode{}
	}
	if cfg.ControllerMode == "" {
		cfg.ControllerMode = "auto"
	}
	if cfg.PearPort == 0 {
		cfg.PearPort = 9863
	}
	if cfg.Monitor.BaseIntervalMs == 0 {
		cfg.Monitor.BaseIntervalMs = 500
	}
	if cfg.Monitor.HysteresisCount == 0 {
		cfg.Monitor.HysteresisCount = 3
	}
	if cfg.Monitor.WatchdogTimeoutSec == 0 {
		cfg.Monitor.WatchdogTimeoutSec = 30
	}
	if cfg.Monitor.ModeDetectThreshold == 0 {
		cfg.Monitor.ModeDetectThreshold = 10
	}
	return cfg
}
