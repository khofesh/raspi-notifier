# raspi-notifier

## Project Overview
A lightweight Go service that polls Gmail and Google Calendar, then sends Telegram notifications for new emails and upcoming meetings. Designed to run on a Raspberry Pi 4 with DietPi.

## Architecture
- Single Go binary, no external runtime dependencies
- Polls on a configurable interval (default: 5 minutes)
- Persists state (Gmail historyId, notified event IDs) to a JSON file to avoid duplicate notifications
- Runs as a systemd service under the `fahmi` user

## File Structure
```
main.go       — config loading, main poll loop
auth.go       — Google OAuth2 flow (Desktop app credentials)
gmail.go      — Gmail API: detect new INBOX messages via History API
calendar.go   — Calendar API: detect events starting within warning window
telegram.go   — Telegram Bot API sender
state.go      — load/save state.json
config.yaml   — user configuration
notifier.service — systemd unit file
```

## Key Design Decisions
- Uses Gmail History API (not full message list) to efficiently detect only new messages since last check
- Calendar events are deduplicated by event ID stored in state file
- State file is written after every check cycle
- OAuth token is stored in a separate file from credentials; first run triggers browser OAuth flow

## Configuration
All config lives in `config.yaml` (path overridable via `NOTIFIER_CONFIG` env var):
- `check_interval_minutes` — polling interval
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
