package main

import (
	"encoding/json"
	"log"
	"os"
)

type State struct {
	GmailHistoryID string          `json:"gmail_history_id"`
	NotifiedEvents map[string]bool `json:"notified_events"`
}

func loadState(path string) *State {
	state := &State{
		NotifiedEvents: make(map[string]bool),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return state
	}
	if err := json.Unmarshal(data, state); err != nil {
		log.Printf("parse state: %v", err)
	}
	if state.NotifiedEvents == nil {
		state.NotifiedEvents = make(map[string]bool)
	}
	return state
}

func saveState(path string, state *State) {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Printf("marshal state: %v", err)
		return
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		log.Printf("write state: %v", err)
	}
}
