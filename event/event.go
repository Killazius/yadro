package event

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Type int

const (
	Register Type = iota + 1
	SetStartTime
	OnStartLine
	Started
	OnFiringRange
	Hit
	LeftFiringRange
	EnteredPenaltyLaps
	LeftPenaltyLaps
	EndMainLap
	CantContinue
)

type Event struct {
	Time         time.Time
	ID           Type
	CompetitorID int
	ExtraParams  string
}

func parse(line string) (*Event, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid line format: %s", line)
	}

	timestamp := strings.Trim(parts[0], "[]")
	eventTime, err := time.Parse("15:04:05.000", timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse event ID: %v", err)
	}
	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse competitor ID: %v", err)
	}

	var extraParams string
	if len(parts) > 3 {
		extraParams = strings.Join(parts[3:], " ")
	}

	return &Event{
		Time:         eventTime,
		ID:           Type(id),
		CompetitorID: competitorID,
		ExtraParams:  extraParams,
	}, nil
}

func Load(r io.Reader) ([]*Event, error) {
	scanner := bufio.NewScanner(r)
	events := make([]*Event, 0)
	for scanner.Scan() {
		line := scanner.Text()
		event, err := parse(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line: %v", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read events: %v", err)
	}
	return events, nil
}
