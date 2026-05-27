package config

import "image"

type Config struct {
	Version        int                    `json:"version"`
	States         map[string]StateConfig `json:"states"`
	GameModes      map[string]GameMode    `json:"game_modes"`
	ControllerMode string                 `json:"controller_mode"`
	PearPort       int                    `json:"pear_port"`
	Monitor        MonitorConfig          `json:"monitor"`
	Discord        DiscordConfig          `json:"discord"`
}

type MonitorConfig struct {
	BaseIntervalMs      int `json:"base_interval_ms"`
	HysteresisCount     int `json:"hysteresis_count"`
	WatchdogTimeoutSec  int `json:"watchdog_timeout_sec"`
	ModeDetectThreshold int `json:"mode_detect_threshold"`
	DisplayIndex        int `json:"display_index"`
}

type StateConfig struct {
	Name        string             `json:"name"`
	Action      Action             `json:"action"`
	IntervalMs  int                `json:"interval_ms"`
	Transitions []TransitionConfig `json:"transitions"`
	BuyTimer    *BuyTimerConfig    `json:"buy_timer,omitempty"`
}

type TransitionConfig struct {
	ID         string          `json:"id"`
	ToState    string          `json:"to_state"`
	References []string        `json:"references"`
	Rect       image.Rectangle `json:"rect"`
	Threshold  int             `json:"threshold"`
	Action     Action          `json:"action"`
}

type BuyTimerConfig struct {
	TimerRect     image.Rectangle   `json:"timer_rect"`
	DigitRects    []image.Rectangle `json:"digit_rects"`
	DigitImages   [10]string        `json:"digit_images"`
	PlayBeforeSec int               `json:"play_before_sec"`
}

type GameMode struct {
	ID              string            `json:"id"`
	DisplayName     string            `json:"display_name"`
	DetectImage     string            `json:"detect_image"`
	DetectRect      image.Rectangle   `json:"detect_rect"`
	ActionOverrides map[string]Action `json:"action_overrides"`
	SkipStates      []string          `json:"skip_states"`
}

type Action string

const (
	ActionPlay  Action = "play"
	ActionPause Action = "pause"
	ActionNone  Action = "none"
)

type DiscordConfig struct {
	Enabled   bool           `json:"enabled"`
	AppID     string         `json:"app_id"`
	RiotID    string         `json:"riot_id"`
	CustomBtn *DiscordButton `json:"custom_button,omitempty"`
}

type DiscordButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}
