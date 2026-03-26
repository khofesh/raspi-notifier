package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Telegram struct {
	token  string
	chatID string
}

func NewTelegram(token, chatID string) *Telegram {
	return &Telegram{token: token, chatID: chatID}
}

func (t *Telegram) Send(text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)

	body, _ := json.Marshal(map[string]string{
		"chat_id": t.chatID,
		"text":    text,
	})

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("telegram send error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("telegram non-200: %s", resp.Status)
	}
}
