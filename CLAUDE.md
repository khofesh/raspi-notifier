# raspi-notifier

## Project Overview
A lightweight Go service that checks Gmail and Google Calendar, then sends Telegram notifications for new emails and upcoming meetings. Designed to run on a Raspberry Pi 4 with DietPi.

## Architecture
- Single Go binary, no external runtime dependencies
- Runs once per invocation (oneshot) — scheduling handled by systemd timer (every 1 minute)
- Persists state (Gmail historyId, notified event IDs) to a JSON file to avoid duplicate notifications
- Runs as a systemd service + timer under the `fahmi` user

## File Structure
```
main.go          — config loading, single check run and exit
auth.go          — Google OAuth2 flow (Desktop app credentials)
gmail.go         — Gmail API: detect new INBOX messages via History API
calendar.go      — Calendar API: detect events starting within warning window
telegram.go      — Telegram Bot API sender
state.go         — load/save state.json
config.yaml      — user configuration
notifier.service — systemd oneshot service unit
notifier.timer   — systemd timer unit (runs every 1 minute)
```

## Key Design Decisions
- Uses systemd `.service` + `.timer` instead of an internal ticker loop — cleaner process lifecycle, `Persistent=true` handles missed runs on reboot
- Uses Gmail History API (not full message list) to efficiently detect only new messages since last check
- Calendar events are deduplicated by event ID stored in state file
- State file is written after every check cycle
- OAuth token is stored in a separate file from credentials; first run triggers browser OAuth flow

## Configuration
All config lives in `config.yaml` (path overridable via `NOTIFIER_CONFIG` env var):
- `meeting_warning_minutes` — how far ahead to warn about meetings
- `telegram_token` / `telegram_chat_id` — Telegram bot credentials
- `state_file` / `credentials_file` / `token_file` — file paths

## Build & Run
```bash
go build -o notifier .
./notifier
```

## Dependencies
- `golang.org/x/oauth2` — Google OAuth2
- `google.golang.org/api` — Gmail + Calendar APIs
- `gopkg.in/yaml.v3` — config parsing
