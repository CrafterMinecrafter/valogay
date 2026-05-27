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
import { Play, Pause, Toggle, GetStatus, GetSongInfo, GetStateInfo, Start, Stop, SetMode, SaveDiscordConfig, SaveMonitorConfig, RecorderStart, RecorderStop, RecorderStatus, RecorderLastFrameImage, GetRegions, SaveRegions, DeleteRegion } from '../wailsjs/go/main/App.js';

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

// Recorder
document.getElementById('btn-recorder-save').addEventListener('click', async () => {
    const enabled = document.getElementById('recorder-enabled').checked;
    if (enabled) {
        await RecorderStart().catch(console.error);
    } else {
        await RecorderStop().catch(console.error);
    }
    await updateRecorderStatus();
});

async function updateRecorderStatus() {
    const rs = await RecorderStatus().catch(() => null);
    if (rs) {
        const el = document.getElementById('recorder-status');
        const chk = document.getElementById('recorder-enabled');
        el.textContent = rs.enabled ? '⏺ Идёт запись' : '⏹ Остановлен';
        chk.checked = rs.enabled;
        document.getElementById('recorder-interval').value = rs.interval_sec;
        document.getElementById('recorder-output').value = rs.output_dir;
    }
}

setInterval(updateRecorderStatus, 3000);
updateRecorderStatus();

// ============================================
// SCREEN EDITOR
// ============================================
(function() {
    const canvas = document.getElementById('screen-canvas');
    const ctx = canvas.getContext('2d', { willReadFrequently: true });
    const wrap = document.getElementById('screen-canvas-wrap');
    const placeholder = document.getElementById('screen-placeholder');
    const infoEl = document.getElementById('screen-info');
    const regionsList = document.getElementById('screen-regions-list');

    let img = null;
    let drawable = null;
    let scale = 1;
    let offsetX = 0, offsetY = 0;
    let mode = 'select'; // 'select' | 'pan'
    let toolSelect = document.getElementById('screen-tool-select');
    let toolPan = document.getElementById('screen-tool-pan');
    let resetBtn = document.getElementById('screen-reset-view');
    let addRegionBtn = document.getElementById('screen-add-region');
    let copyCoordsBtn = document.getElementById('screen-copy-coords');

    let regions = [];
    let selectedRegionId = null;
    let drawing = false;
    let drawStart = null, drawCur = null;
    let dragging = false;
    let dragStartPt = null;

    // --- Load regions ---
    async function loadRegions() {
        regions = await GetRegions().catch(() => []) || [];
        renderRegionsList();
    }

    // --- Save regions ---
    async function saveRegions() {
        await SaveRegions(regions).catch(console.error);
    }

    // --- Render regions list sidebar ---
    function renderRegionsList() {
        regionsList.innerHTML = '';
        regions.forEach(r => {
            const div = document.createElement('div');
            div.className = 'region-item' + (r.id === selectedRegionId ? ' selected' : '');
            const cx = Math.round(r.x + r.w / 2);
            const cy = Math.round(r.y + r.h / 2);
            let pixelHex = '—';
            if (img && r.w > 0 && r.h > 0) {
                const c = getPixelColor(cx, cy);
                pixelHex = rgbToHex(c.r, c.g, c.b);
            }
            div.innerHTML = `
                <div class="region-item-header">
                    <span class="region-item-name">${r.name}</span>
                    <button class="region-item-delete" data-id="${r.id}">&times;</button>
                </div>
                <div class="region-item-coords">x=${r.x} y=${r.y} ${r.w}x${r.h}</div>
                <div class="region-item-pixel">
                    <div class="region-pixel-preview" style="background:${pixelHex}"></div>
                    <span class="region-pixel-hex">${pixelHex}</span>
                </div>
            `;
            div.onclick = (e) => {
                if (e.target.classList.contains('region-item-delete')) return;
                selectedRegionId = r.id;
                renderRegionsList();
                zoomToRegion(r);
            };
            div.querySelector('.region-item-delete').onclick = async (e) => {
                e.stopPropagation();
                await DeleteRegion(r.id).catch(console.error);
                await loadRegions();
            };
            regionsList.appendChild(div);
        });
    }

    function getPixelColor(x, y) {
        if (!img) return {r:0,g:0,b:0};
        const w = img.width, h = img.height;
        if (x < 0 || x >= w || y < 0 || y >= h) return {r:0,g:0,b:0};
        const idx = (y * w + x) * 4;
        return { r: img.data[idx], g: img.data[idx+1], b: img.data[idx+2] };
    }

    function rgbToHex(r, g, b) {
        return '#' + [r,g,b].map(x => x.toString(16).padStart(2,'0')).join('');
    }

    // --- Zoom to region ---
    function zoomToRegion(r) {
        const cw = wrap.clientWidth, ch = wrap.clientHeight;
        scale = Math.min(cw / r.w, ch / r.h) * 0.8;
        if (scale < 0.5) scale = 0.5;
        if (scale > 20) scale = 20;
        offsetX = (cw - r.w * scale) / 2 - r.x * scale;
        offsetY = (ch - r.h * scale) / 2 - r.y * scale;
        render();
    }

    // --- Render ---
    function render() {
        if (!img) {
            placeholder.style.display = 'flex';
            return;
        }
        placeholder.style.display = 'none';
        const w = wrap.clientWidth, h = wrap.clientHeight;
        canvas.width = w;
        canvas.height = h;
        ctx.fillStyle = '#1c2333';
        ctx.fillRect(0, 0, w, h);
        ctx.save();
        ctx.translate(offsetX, offsetY);
        ctx.scale(scale, scale);
        if (drawable) {
            ctx.drawImage(drawable, 0, 0);
        }
        ctx.restore();
        // Draw regions
        ctx.save();
        ctx.translate(offsetX, offsetY);
        ctx.scale(scale, scale);
        regions.forEach(r => {
            ctx.strokeStyle = r.id === selectedRegionId ? '#58a6ff' : 'rgba(88,166,255,0.5)';
            ctx.lineWidth = r.id === selectedRegionId ? 3 / scale : 1 / scale;
            ctx.setLineDash([4/scale, 4/scale]);
            ctx.strokeRect(r.x, r.y, r.w, r.h);
            if (r.id === selectedRegionId) {
                ctx.fillStyle = 'rgba(88,166,255,0.1)';
                ctx.fillRect(r.x, r.y, r.w, r.h);
            }
        });
        ctx.restore();
        // Drawing rect
        if (drawing && drawStart && drawCur) {
            const sx = Math.min(drawStart.x, drawCur.x);
            const sy = Math.min(drawStart.y, drawCur.y);
            const sw = Math.abs(drawCur.x - drawStart.x);
            const sh = Math.abs(drawCur.y - drawStart.y);
            ctx.save();
            ctx.translate(offsetX, offsetY);
            ctx.scale(scale, scale);
            ctx.strokeStyle = '#58a6ff';
            ctx.lineWidth = 2 / scale;
            ctx.setLineDash([4/scale, 4/scale]);
            ctx.strokeRect(sx, sy, sw, sh);
            ctx.restore();
        }
    }

    // --- Load image ---
    async function loadFrame() {
        const b64 = await RecorderLastFrameImage().catch(() => '');
        if (!b64) return;
        const dataUrl = 'data:image/png;base64,' + b64;
        const im = new Image();
        im.src = dataUrl;
        await new Promise(r => im.onload = r);
        const c = document.createElement('canvas');
        c.width = im.width; c.height = im.height;
        const cx = c.getContext('2d');
        cx.drawImage(im, 0, 0);
        img = cx.getImageData(0, 0, im.width, im.height);
        const bitmap = await createImageBitmap(img);
        drawable = bitmap;
        renderRegionsList();
        render();
    }

    // --- Coords helpers ---
    function screenToImage(sx, sy) {
        return {
            x: Math.floor((sx - offsetX) / scale),
            y: Math.floor((sy - offsetY) / scale)
        };
    }

    // --- Event handlers ---
    toolSelect.onclick = () => { mode = 'select'; toolSelect.classList.add('active'); toolPan.classList.remove('active'); canvas.style.cursor = 'crosshair'; };
    toolPan.onclick = () => { mode = 'pan'; toolPan.classList.add('active'); toolSelect.classList.remove('active'); canvas.style.cursor = 'grab'; };
    resetBtn.onclick = () => { scale = 1; offsetX = 0; offsetY = 0; render(); };

    addRegionBtn.onclick = async () => {
        const name = 'Region ' + (regions.length + 1);
        const newRegion = {
            id: crypto.randomUUID(),
            name: name,
            x: Math.floor(img ? img.width / 4 : 0),
            y: Math.floor(img ? img.height / 4 : 0),
            w: Math.floor(img ? img.width / 2 : 100),
            h: Math.floor(img ? img.height / 2 : 100)
        };
        regions.push(newRegion);
        selectedRegionId = newRegion.id;
        await saveRegions();
        renderRegionsList();
        render();
    };

    copyCoordsBtn.onclick = () => {
        const txt = JSON.stringify(regions.map(r => ({ name: r.name, x: r.x, y: r.y, w: r.w, h: r.h })), null, 2);
        navigator.clipboard.writeText(txt);
        copyCoordsBtn.textContent = 'Скопировано!';
        setTimeout(() => copyCoordsBtn.textContent = 'Копировать coords', 1500);
    };

    canvas.onmousedown = (e) => {
        const rect = canvas.getBoundingClientRect();
        const mx = e.clientX - rect.left;
        const my = e.clientY - rect.top;
        if (mode === 'pan') {
            dragging = true;
            dragStartPt = { x: mx, y: my };
            canvas.style.cursor = 'grabbing';
        } else {
            drawing = true;
            drawStart = screenToImage(mx, my);
            drawCur = { ...drawStart };
        }
    };

    canvas.onmousemove = (e) => {
        const rect = canvas.getBoundingClientRect();
        const mx = e.clientX - rect.left;
        const my = e.clientY - rect.top;
        const pos = screenToImage(mx, my);
        let info = `${pos.x}, ${pos.y}`;
        if (img) {
            const c = getPixelColor(pos.x, pos.y);
            info += ` | ${rgbToHex(c.r,c.g,c.b)}`;
        }
        if (mode === 'pan' && dragging) {
            offsetX += mx - dragStartPt.x;
            offsetY += my - dragStartPt.y;
            dragStartPt = { x: mx, y: my };
            render();
        } else if (drawing) {
            drawCur = pos;
            render();
        }
        infoEl.textContent = info;
    };

    canvas.onmouseup = async (e) => {
        if (mode === 'pan' && dragging) {
            dragging = false;
            canvas.style.cursor = 'grab';
        }
        if (drawing && drawStart && drawCur) {
            const sx = Math.min(drawStart.x, drawCur.x);
            const sy = Math.min(drawStart.y, drawCur.y);
            const sw = Math.abs(drawCur.x - drawStart.x);
            const sh = Math.abs(drawCur.y - drawStart.y);
            if (sw > 5 && sh > 5) {
                const name = 'Region ' + (regions.length + 1);
                const newRegion = {
                    id: crypto.randomUUID(),
                    name: name,
                    x: sx, y: sy, w: sw, h: sh
                };
                regions.push(newRegion);
                selectedRegionId = newRegion.id;
                await saveRegions();
                renderRegionsList();
            }
            drawing = false;
            drawStart = null;
            drawCur = null;
            render();
        }
    };

    canvas.onwheel = (e) => {
        e.preventDefault();
        const rect = canvas.getBoundingClientRect();
        const mx = e.clientX - rect.left;
        const my = e.clientY - rect.top;
        const oldScale = scale;
        scale *= e.deltaY < 0 ? 1.2 : 0.8;
        scale = Math.max(0.2, Math.min(30, scale));
        const zoomFactor = scale / oldScale;
        offsetX = mx - (mx - offsetX) * zoomFactor;
        offsetY = my - (my - offsetY) * zoomFactor;
        render();
    };

    // --- Init ---
    toolSelect.classList.add('active');
    canvas.style.cursor = 'crosshair';

    // Load frame + regions when tab shown
    const observer = new MutationObserver(() => {
        const tab = document.getElementById('screen');
        if (!tab.classList.contains('hidden')) {
            loadFrame();
            loadRegions();
        }
    });
    observer.observe(document.getElementById('app'), { attributes: true, subtree: true });

    window._screenEditor = { loadFrame, render, zoomToRegion: (r) => { selectedRegionId = r?.id; zoomToRegion(r); render(); } };
})();

// Polling for new frames (when on screen tab)
setInterval(() => {
    const tab = document.getElementById('screen');
    if (!tab || tab.classList.contains('hidden')) return;
    if (window._screenEditor && window._screenEditor.loadFrame) {
        window._screenEditor.loadFrame().catch(() => {});
    }
}, 3000);