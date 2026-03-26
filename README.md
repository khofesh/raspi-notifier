# raspi-notifier

A lightweight Go service that monitors Gmail and Google Calendar and sends Telegram messages when you receive a new email or have a meeting coming up.

Designed to run on a Raspberry Pi 4 with DietPi as a systemd oneshot service triggered by a timer every minute.

## Features

- New email in Gmail inbox → Telegram message with sender and subject
- Upcoming calendar meeting → Telegram message with time remaining and title
- Configurable meeting warning window
- Minimal resource usage — single binary, exits after each run, no long-running process

## Prerequisites

- Go 1.26.1+
- A [Google Cloud project](https://console.cloud.google.com/) with **Gmail API** and **Google Calendar API** enabled
- OAuth2 credentials (type: **Desktop app**) downloaded as `credentials.json`
- A Telegram bot token from [@BotFather](https://t.me/BotFather)
- Your Telegram chat ID (get it from [@userinfobot](https://t.me/userinfobot))

## Setup

### 1. Google Cloud credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/) and select or create a project
2. **APIs & Services** → **Library** → search and enable **Gmail API**
3. **APIs & Services** → **Library** → search and enable **Google Calendar API**
4. **APIs & Services** → **OAuth consent screen** → set User Type to **External** → fill in app name and your email → **Save**
5. On the same OAuth consent screen page → **Test users** → **Add users** → add your Gmail address
6. **APIs & Services** → **Credentials** → **Create Credentials** → **OAuth client ID** → Application type: **Desktop app** → **Create**
7. Click the download button (↓) next to the credential you just created
8. Rename and place the file:
   ```bash
   cp ~/Downloads/client_secret_*.json ~/.notifier_credentials.json
   chmod 600 ~/.notifier_credentials.json
   ```

### 2. Configure

```bash
mkdir -p ~/.notifier
cp config.yaml ~/.notifier/config.yaml
```

Edit `~/.notifier/config.yaml`:

```yaml
meeting_warning_minutes: 15
telegram_token: "YOUR_BOT_TOKEN"
telegram_chat_id: "YOUR_CHAT_ID"
state_file: "/home/fahmi/.notifier/state.json"
credentials_file: "/home/fahmi/.notifier_credentials.json"
token_file: "/home/fahmi/.notifier_token.json"
```

### 3. Build

```bash
go build -o notifier .
mkdir -p ~/bin
cp notifier ~/bin/notifier
```

### 4. First run (OAuth)

Run once manually on the Pi — it will print a Google OAuth URL. Open it in a browser on your PC, authorize, and paste the code back into the terminal:

```bash
NOTIFIER_CONFIG=~/.notifier/config.yaml ~/bin/notifier
```

This creates the token file. Subsequent runs are fully automatic.

### 5. Install as systemd service + timer

```bash
sudo cp notifier.service /etc/systemd/system/
sudo cp notifier.timer /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now notifier.timer
```

Check the timer status:

```bash
systemctl status notifier.timer
systemctl list-timers notifier.timer
```

Check logs:

```bash
journalctl -u notifier.service -f
```

Trigger a manual run:

```bash
systemctl start notifier.service
```

## Configuration Reference

| Key | Description |
|---|---|
| `meeting_warning_minutes` | Warn about meetings starting within this window |
| `telegram_token` | Telegram bot token from @BotFather |
| `telegram_chat_id` | Your Telegram user/chat ID |
| `state_file` | Path to persist state between runs |
| `credentials_file` | Path to Google OAuth credentials JSON |
| `token_file` | Path to store the OAuth token after first auth |
