// tab switching
document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.add('hidden'));
        btn.classList.add('active');
        const tabId = btn.dataset.tab;
        document.getElementById(tabId).classList.remove('hidden');
    });
});

// Import from wailsjs generated bindings
import { Play, Pause, Toggle, GetStatus, GetSongInfo, GetStateInfo, Start, Stop, SetMode, SaveDiscordConfig, SaveMonitorConfig } from '../wailsjs/go/main/App.js';

// Controller buttons
document.getElementById('btn-play').addEventListener('click', () => {
    Play().catch(console.error);
});

document.getElementById('btn-pause').addEventListener('click', () => {
    Pause().catch(console.error);
});

document.getElementById('btn-toggle').addEventListener('click', () => {
    Toggle().catch(console.error);
});

// Update status periodically
async function updateStatus() {
    const status = await GetStatus().catch(() => null);
    if (status) {
        document.getElementById('winkey-status').textContent = 'WinKey: ' + (status.WinKey ? 'ON' : 'OFF');
        document.getElementById('pear-status').textContent = 'Pear: ' + (status.Pear ? 'ON' : 'OFF');
    }

    const info = await GetSongInfo().catch(() => null);
    if (info) {
        document.getElementById('song-title').textContent = info.title || 'Нет трека';
        document.getElementById('song-artist').textContent = info.artist || '—';
    }

    const si = await GetStateInfo().catch(() => null);
    if (si) {
        document.getElementById('status-bar').textContent = '● ' + si.state + ' | ⚔️ ' + si.mode + ' | ' + si.action + ' | Раунд ' + si.round;
    }
}

setInterval(updateStatus, 1000);
updateStatus();

// Apply/Stop buttons
document.getElementById('btn-apply').addEventListener('click', () => {
    Start().catch(console.error);
});

document.getElementById('btn-stop').addEventListener('click', () => {
    Stop().catch(console.error);
});

// Save settings
document.getElementById('btn-save-settings').addEventListener('click', async () => {
    const interval = parseInt(document.getElementById('cfg-interval').value) || 500;
    const hysteresis = parseInt(document.getElementById('cfg-hysteresis').value) || 3;
    const watchdog = parseInt(document.getElementById('cfg-watchdog').value) || 30;
    await SaveMonitorConfig(interval, hysteresis, watchdog, 10);
});

// Save Discord config
document.getElementById('btn-save-discord').addEventListener('click', async () => {
    const enabled = document.getElementById('discord-enabled').checked;
    const appID = document.getElementById('discord-appid').value;
    const riotID = document.getElementById('discord-riotid').value;
    const btnLabel = document.getElementById('discord-btn-label').value;
    const btnURL = document.getElementById('discord-btn-url').value;
    await SaveDiscordConfig(enabled, appID, riotID, btnLabel, btnURL);
});

// Mode select
const modeSelect = document.getElementById('controller-mode');
if (modeSelect) {
    modeSelect.addEventListener('change', () => {
        SetMode(modeSelect.value).catch(console.error);
    });
}