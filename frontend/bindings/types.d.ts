// Wails TypeScript bindings for VPMC

export interface Status {
	WinKey: boolean;
	Pear: boolean;
}

export interface SongInfo {
	title: string;
	artist: string;
	isPaused: boolean;
}

export interface DiscordConfig {
	enabled: boolean;
	app_id: string;
	riot_id: string;
}

export interface MonitorConfig {
	base_interval_ms: number;
	hysteresis_count: number;
	watchdog_timeout_sec: number;
	mode_detect_threshold: number;
}
