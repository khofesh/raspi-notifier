package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"
)

type EventSummary struct {
	ID       string
	Summary  string
	StartsIn string
}

func checkCalendar(ctx context.Context, svc *calendar.Service, warningMins int, notified map[string]bool) ([]EventSummary, error) {
	now := time.Now()
	windowEnd := now.Add(time.Duration(warningMins) * time.Minute)

	resp, err := svc.Events.List("primary").
		TimeMin(now.Format(time.RFC3339)).
		TimeMax(windowEnd.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	var upcoming []EventSummary
	for _, item := range resp.Items {
		if notified[item.Id] {
			continue
		}

		startStr := item.Start.DateTime
		if startStr == "" {
			startStr = item.Start.Date
		}
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			continue
		}

		startsIn := time.Until(start).Round(time.Minute)
		upcoming = append(upcoming, EventSummary{
			ID:       item.Id,
			Summary:  item.Summary,
			StartsIn: startsIn.String(),
		})
	}

	return upcoming, nil
}
