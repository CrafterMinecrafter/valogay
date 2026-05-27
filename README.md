# VPMC — Valorant Phase Music Controller

VPMC — фоновое приложение для Windows 10/11, которое управляет воспроизведением музыки во время матчей Valorant на основе анализа экрана (без внедрения в память игры).

## Что делает проект

- Захватывает области экрана (DXGI Desktop Duplication).
- Сравнивает кадры с эталонами через pHash.
- Ведет конечный автомат состояний матча (FSM).
- Определяет режим игры (competitive/unrated/spike_rush/deathmatch).
- Отправляет Play/Pause через WinKey или Pear API.
- Обновляет Discord Rich Presence.
- Предоставляет GUI (Fyne) для настройки переходов/областей и статуса.

## Архитектура каталогов

- `config/` — типы конфигурации, load/save, миграции, zip-профили.
- `fsm/` — FSM, режимы, таймер BUY, тесты.
- `vision/` — capturer, compare (pHash cache), normalize, autodetect, OCR timer.
- `controller/` — интерфейс контроллера, WinKey, Pear, manager.
- `monitor/` — цикл мониторинга, hysteresis, watchdog.
- `discord/` — PresenceManager + тесты.
- `gui/` — приложение Fyne, табы, трей, виджеты.
- `assets/` — иконки/ресурсы.

## Используемые библиотеки (по ТЗ)

- GUI: `fyne.io/fyne/v2`
- Screen capture: `github.com/kirides/go-d3d`
- pHash: `github.com/corona10/goimagehash`
- Image processing: `github.com/disintegration/imaging`
- Win API: `golang.org/x/sys/windows`
- Discord RPC: `github.com/axrona/go-discordrpc/client`
- Tests: `github.com/stretchr/testify`

## Текущий статус реализации

> Ниже фиксируется текущее состояние репозитория на этой ветке.

### Реализовано

- Базовые структуры конфигурации, load/save/default/migrate, ZIP профили.
- FSM с состоянием/режимом/round counter/callbacks.
- Controller manager + Pear controller + WinKey заглушки/каркас.
- Monitor loop c hysteresis/watchdog и mode-detect веткой.
- Vision comparer cache + normalize + Capturer интерфейсы и DXGI каркас.
- Discord presence manager каркас с правилами кнопок.
- Базовая структура `gui/` и `widgets/`.
- `vpmc.exe.manifest` с `PerMonitorV2`.

### В процессе / TODO до полного соответствия ТЗ

- Полноценная DXGI реализация на `kirides/go-d3d` (dirty rect feed + буферы GPU).
- Полный pHash pipeline именно через `goimagehash` + `imaging`.
- Реальный WinKey через `user32 keybd_event` и точное поведение `lastKnownAction`.
- Полноценный Discord transport через `axrona/go-discordrpc/client` и reconnection loop.
- Полный GUI на Fyne: табы, region editor (zoom/pan/live test), трей и статусная строка.
- Первый запуск wizard, crash recovery (`last_state.json`), hotkey `Ctrl+Alt+M`.
- Дополнительные unit/integration тесты: Pear mock server, Discord mock client, monitor logic.

## Сборка

### Требования

- Go 1.21+
- Windows 10/11 для production-запуска

### Локальная проверка

```bash
go test ./...
```

### Сборка бинарника (Windows)

```bash
go build -o vpmc.exe ./
```

Убедитесь, что рядом с `vpmc.exe` присутствует `vpmc.exe.manifest`.

## Конфиг

По умолчанию приложение ищет `config.json` в рабочей директории.
Если файла нет — создается/используется default-конфиг в памяти.

## Безопасность

VPMC не читает память Valorant и не внедряется в процесс игры; анализ идет только по изображениям экрана.
