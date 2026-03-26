package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func newGoogleServices(ctx context.Context, cfg *Config) (*gmail.Service, *calendar.Service, error) {
	credData, err := os.ReadFile(cfg.CredentialsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("read credentials: %w", err)
	}

	oauthCfg, err := google.ConfigFromJSON(credData,
		gmail.GmailReadonlyScope,
		calendar.CalendarReadonlyScope,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("parse credentials: %w", err)
	}

	token, err := tokenFromFile(cfg.TokenFile)
	if err != nil {
		token = getTokenFromWeb(oauthCfg)
		saveToken(cfg.TokenFile, token)
	}

	httpClient := oauthCfg.Client(ctx, token)

	gmailSvc, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, fmt.Errorf("gmail service: %w", err)
	}

	calSvc, err := calendar.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, fmt.Errorf("calendar service: %w", err)
	}

	return gmailSvc, calSvc, nil
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	if err := json.NewDecoder(f).Decode(&tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func getTokenFromWeb(cfg *oauth2.Config) *oauth2.Token {
	authURL := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Open this URL in your browser and paste the code:\n%v\n\nCode: ", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("read auth code: %v", err)
	}

	tok, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("exchange token: %v", err)
	}
	return tok
}

func saveToken(path string, tok *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("save token: %v", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)
}

