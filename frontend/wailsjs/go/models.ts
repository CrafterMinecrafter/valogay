export namespace config {
	
	export class BuyTimerConfig {
	    timer_rect: image.Rectangle;
	    digit_rects: image.Rectangle[];
	    digit_images: string[];
	    play_before_sec: number;
	
	    static createFrom(source: any = {}) {
	        return new BuyTimerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timer_rect = this.convertValues(source["timer_rect"], image.Rectangle);
	        this.digit_rects = this.convertValues(source["digit_rects"], image.Rectangle);
	        this.digit_images = source["digit_images"];
	        this.play_before_sec = source["play_before_sec"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Region {
	    id: string;
	    name: string;
	    x: number;
	    y: number;
	    w: number;
	    h: number;
	
	    static createFrom(source: any = {}) {
	        return new Region(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.w = source["w"];
	        this.h = source["h"];
	    }
	}
	export class DiscordButton {
	    label: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new DiscordButton(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.url = source["url"];
	    }
	}
	export class DiscordConfig {
	    enabled: boolean;
	    app_id: string;
	    riot_id: string;
	    custom_button?: DiscordButton;
	
	    static createFrom(source: any = {}) {
	        return new DiscordConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.app_id = source["app_id"];
	        this.riot_id = source["riot_id"];
	        this.custom_button = this.convertValues(source["custom_button"], DiscordButton);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RecorderConfig {
	    enabled: boolean;
	    interval_sec: number;
	    output_dir: string;
	
	    static createFrom(source: any = {}) {
	        return new RecorderConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.interval_sec = source["interval_sec"];
	        this.output_dir = source["output_dir"];
	    }
	}
	export class MonitorConfig {
	    base_interval_ms: number;
	    hysteresis_count: number;
	    watchdog_timeout_sec: number;
	    mode_detect_threshold: number;
	    display_index: number;
	
	    static createFrom(source: any = {}) {
	        return new MonitorConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.base_interval_ms = source["base_interval_ms"];
	        this.hysteresis_count = source["hysteresis_count"];
	        this.watchdog_timeout_sec = source["watchdog_timeout_sec"];
	        this.mode_detect_threshold = source["mode_detect_threshold"];
	        this.display_index = source["display_index"];
	    }
	}
	export class GameMode {
	    id: string;
	    display_name: string;
	    detect_image: string;
	    detect_rect: image.Rectangle;
	    action_overrides: Record<string, string>;
	    skip_states: string[];
	
	    static createFrom(source: any = {}) {
	        return new GameMode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.display_name = source["display_name"];
	        this.detect_image = source["detect_image"];
	        this.detect_rect = this.convertValues(source["detect_rect"], image.Rectangle);
	        this.action_overrides = source["action_overrides"];
	        this.skip_states = source["skip_states"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TransitionConfig {
	    id: string;
	    to_state: string;
	    references: string[];
	    rect: image.Rectangle;
	    threshold: number;
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new TransitionConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.to_state = source["to_state"];
	        this.references = source["references"];
	        this.rect = this.convertValues(source["rect"], image.Rectangle);
	        this.threshold = source["threshold"];
	        this.action = source["action"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StateConfig {
	    name: string;
	    action: string;
	    interval_ms: number;
	    transitions: TransitionConfig[];
	    buy_timer?: BuyTimerConfig;
	
	    static createFrom(source: any = {}) {
	        return new StateConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.action = source["action"];
	        this.interval_ms = source["interval_ms"];
	        this.transitions = this.convertValues(source["transitions"], TransitionConfig);
	        this.buy_timer = this.convertValues(source["buy_timer"], BuyTimerConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Config {
	    version: number;
	    states: Record<string, StateConfig>;
	    game_modes: Record<string, GameMode>;
	    controller_mode: string;
	    pear_port: number;
	    pear_auth_id?: string;
	    pear_token?: string;
	    monitor: MonitorConfig;
	    recorder: RecorderConfig;
	    discord: DiscordConfig;
	    regions: Region[];
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.states = this.convertValues(source["states"], StateConfig, true);
	        this.game_modes = this.convertValues(source["game_modes"], GameMode, true);
	        this.controller_mode = source["controller_mode"];
	        this.pear_port = source["pear_port"];
	        this.pear_auth_id = source["pear_auth_id"];
	        this.pear_token = source["pear_token"];
	        this.monitor = this.convertValues(source["monitor"], MonitorConfig);
	        this.recorder = this.convertValues(source["recorder"], RecorderConfig);
	        this.discord = this.convertValues(source["discord"], DiscordConfig);
	        this.regions = this.convertValues(source["regions"], Region);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	

}

export namespace image {
	
	export class Point {
	    X: number;
	    Y: number;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	    }
	}
	export class Rectangle {
	    Min: Point;
	    Max: Point;
	
	    static createFrom(source: any = {}) {
	        return new Rectangle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Min = this.convertValues(source["Min"], Point);
	        this.Max = this.convertValues(source["Max"], Point);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class SongInfo {
	    title: string;
	    artist: string;
	    isPaused: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SongInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.isPaused = source["isPaused"];
	    }
	}
	export class Status {
	    WinKey: boolean;
	    Pear: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WinKey = source["WinKey"];
	        this.Pear = source["Pear"];
	    }
	}

}

