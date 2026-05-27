// This file is auto-generated. Do not edit manually.

import { invoke } from '@wailsio/runtime';

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
}

export interface MonitorConfig {
	base_interval_ms: number;
	hysteresis_count: number;
	watchdog_timeout_sec: number;
	mode_detect_threshold: number;
}

export interface StateConfig {
	name: string;
	action: string;
	interval_ms: number;
}

export interface GameMode {
	id: string;
	display_name: string;
}

export function Play(): Promise<void> {
	return invoke('app.play');
}

export function Pause(): Promise<void> {
	return invoke('app.pause');
}

export function Toggle(): Promise<void> {
	return invoke('app.toggle');
}

export function GetStatus(): Promise<Status> {
	return invoke('app.getStatus');
}

export function GetSongInfo(): Promise<SongInfo | null> {
	return invoke('app.getSongInfo');
}

export function GetStateInfo(): Promise<Record<string, any>> {
	return invoke('app.getStateInfo');
}

export function GetConfig(): Promise<any> {
	return invoke('app.getConfig');
}

export function SaveConfig(cfg: any): Promise<void> {
	return invoke('app.saveConfig');
}

export function SetMode(mode: string): Promise<void> {
	return invoke('app.setMode', { mode });
}

export function SetPearPort(port: number): Promise<void> {
	return invoke('app.setPearPort', { port });
}

export function Start(): Promise<void> {
	return invoke('app.start');
}

export function Stop(): Promise<void> {
	return invoke('app.stop');
}

export function GetDiscordConfig(): Promise<DiscordConfig> {
	return invoke('app.getDiscordConfig');
}

export function SaveDiscordConfig(enabled: boolean, appID: string, riotID: string, btnLabel: string, cmdURL: string): Promise<void> {
	return invoke('app.saveDiscordConfig', { enabled, appID, riotID, btnLabel, cmdURL });
}

export function GetMonitorConfig(): Promise<MonitorConfig> {
	return invoke('app.getMonitorConfig');
}

export function SaveMonitorConfig(intervalMs: number, hysteresis: number, watchdogSec: number, modeThreshold: number): Promise<void> {
	return invoke('app.saveMonitorConfig', { intervalMs, hysteresis, watchdogSec, modeThreshold });
}

export function GetStates(): Promise<Record<string, StateConfig>> {
	return invoke('app.getStates');
}

export function GetModes(): Promise<Record<string, GameMode>> {
	return invoke('app.getModes');
}
