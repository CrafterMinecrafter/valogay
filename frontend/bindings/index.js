const { invoke } = window.runtime;

export function Play() {
	return invoke('app.play');
}

export function Pause() {
	return invoke('app.pause');
}

export function Toggle() {
	return invoke('app.toggle');
}

export function GetStatus() {
	return invoke('app.getStatus');
}

export function GetSongInfo() {
	return invoke('app.getSongInfo');
}

export function GetStateInfo() {
	return invoke('app.getStateInfo');
}

export function GetConfig() {
	return invoke('app.getConfig');
}

export function SaveConfig(cfg) {
	return invoke('app.saveConfig', cfg);
}

export function SetMode(mode) {
	return invoke('app.setMode', { mode });
}

export function SetPearPort(port) {
	return invoke('app.setPearPort', { port });
}

export function Start() {
	return invoke('app.start');
}

export function Stop() {
	return invoke('app.stop');
}

export function GetDiscordConfig() {
	return invoke('app.getDiscordConfig');
}

export function SaveDiscordConfig(enabled, appID, riotID, btnLabel, btnURL) {
	return invoke('app.saveDiscordConfig', { enabled, appID, riotID, btnLabel, btnURL });
}

export function GetMonitorConfig() {
	return invoke('app.getMonitorConfig');
}

export function SaveMonitorConfig(intervalMs, hysteresis, watchdogSec, modeThreshold) {
	return invoke('app.saveMonitorConfig', { intervalMs, hysteresis, watchdogSec, modeThreshold });
}

export function GetStates() {
	return invoke('app.getStates');
}

export function GetModes() {
	return invoke('app.getModes');
}