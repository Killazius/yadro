package event

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Type - описание события по ID
type Type int

const (
	Register           Type = iota + 1 // 1: The competitor registered
	SetStartTime                       // 2: The start time was set by a draw
	OnStartLine                        // 3: The competitor is on the start line
	Started                            // 4: The competitor has started
	OnFiringRange                      // 5: The competitor is on the firing range
	Hit                                // 6: The target has been hit
	LeftFiringRange                    // 7: The competitor left the firing range
	EnteredPenaltyLaps                 // 8: The competitor entered the penalty laps
	LeftPenaltyLaps                    // 9: The competitor left the penalty laps
	EndMainLap                         // 10: The competitor ended the main lap
	CantContinue                       // 11: The competitor can't continue
)

// Event - структура для описания события в гонке
type Event struct {
	Time         time.Time // Время возникновения события
	ID           Type      // Тип события (из перечисления Type)
	CompetitorID int       // Идентификатор участника, к которому относится событие
	ExtraParams  string    // Дополнительные параметры в формате строки
}

// parse обрабатывает строку с событием и возвращает Event или ошибку
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

// Load считывает данные из Reader, парсит при помощи parse и возращает слайс из событий
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
