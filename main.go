package main

import (
	"context"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MeetingWarningMins int    `yaml:"meeting_warning_minutes"`
	TelegramToken      string `yaml:"telegram_token"`
	TelegramChatID     string `yaml:"telegram_chat_id"`
	StateFile          string `yaml:"state_file"`
	CredentialsFile    string `yaml:"credentials_file"`
	TokenFile          string `yaml:"token_file"`
	LogFile            string `yaml:"log_file"`
}

func loadConfig() *Config {
	path := "config.yaml"
	if p := os.Getenv("NOTIFIER_CONFIG"); p != "" {
		path = p
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read config: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("parse config: %v", err)
	}
	return &cfg
}

func check(ctx context.Context, cfg *Config) {
	tg := NewTelegram(cfg.TelegramToken, cfg.TelegramChatID)

	gmailSvc, calSvc, err := newGoogleServices(ctx, cfg)
	if err != nil {
		log.Printf("google auth error: %v", err)
		return
	}

	state := loadState(cfg.StateFile)

	// Check Gmail
	messages, newHistoryID, err := checkGmail(ctx, gmailSvc, state.GmailHistoryID)
	if err != nil {
		log.Printf("gmail error: %v", err)
	} else {
		for _, msg := range messages {
			tg.Send("📧 New email from " + msg.From + "\nSubject: " + msg.Subject)
		}
		if newHistoryID != "" {
			state.GmailHistoryID = newHistoryID
		}
	}

	// Check Calendar
	events, err := checkCalendar(ctx, calSvc, cfg.MeetingWarningMins, state.NotifiedEvents)
	if err != nil {
		log.Printf("calendar error: %v", err)
	} else {
		for _, ev := range events {
			tg.Send("📅 Meeting in ~" + ev.StartsIn + ": " + ev.Summary)
			state.NotifiedEvents[ev.ID] = true
		}
	}

	saveState(cfg.StateFile, state)
}

func setupLogger(cfg *Config) {
	if cfg.LogFile == "" {
		return
	}
	f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Printf("warning: could not open log file %s: %v — logging to stderr", cfg.LogFile, err)
		return
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	cfg := loadConfig()
	setupLogger(cfg)
	ctx := context.Background()
	check(ctx, cfg)
}
