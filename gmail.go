package main

import (
	"context"
	"fmt"

	"google.golang.org/api/gmail/v1"
)

type EmailSummary struct {
	From    string
	Subject string
}

func checkGmail(ctx context.Context, svc *gmail.Service, lastHistoryID string) ([]EmailSummary, string, error) {
	var messages []EmailSummary

	if lastHistoryID == "" {
		// First run: get current historyId without fetching all mail
		profile, err := svc.Users.GetProfile("me").Context(ctx).Do()
		if err != nil {
			return nil, "", fmt.Errorf("get profile: %w", err)
		}
		return nil, fmt.Sprintf("%d", profile.HistoryId), nil
	}

	// Fetch history since last known ID
	histResp, err := svc.Users.History.List("me").
		StartHistoryId(mustParseUint64(lastHistoryID)).
		HistoryTypes("messageAdded").
		LabelId("INBOX").
		Context(ctx).
		Do()
	if err != nil {
		return nil, lastHistoryID, fmt.Errorf("list history: %w", err)
	}

	seen := map[string]bool{}
	for _, h := range histResp.History {
		for _, ma := range h.MessagesAdded {
			if seen[ma.Message.Id] {
				continue
			}
			seen[ma.Message.Id] = true

			summary, err := fetchEmailSummary(ctx, svc, ma.Message.Id)
			if err != nil {
				continue
			}
			messages = append(messages, summary)
		}
	}

	newHistoryID := lastHistoryID
	if histResp.HistoryId != 0 {
		newHistoryID = fmt.Sprintf("%d", histResp.HistoryId)
	}

	return messages, newHistoryID, nil
}

func fetchEmailSummary(ctx context.Context, svc *gmail.Service, msgID string) (EmailSummary, error) {
	msg, err := svc.Users.Messages.Get("me", msgID).
		Format("metadata").
		MetadataHeaders("From", "Subject").
		Context(ctx).
		Do()
	if err != nil {
		return EmailSummary{}, err
	}

	var summary EmailSummary
	for _, h := range msg.Payload.Headers {
		switch h.Name {
		case "From":
			summary.From = h.Value
		case "Subject":
			summary.Subject = h.Value
		}
	}
	return summary, nil
}

func mustParseUint64(s string) uint64 {
	var v uint64
	fmt.Sscanf(s, "%d", &v)
	return v
}
