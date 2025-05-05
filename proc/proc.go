package proc

import (
	"fmt"
	"github.com/Killazius/yadro/competitor"
	"github.com/Killazius/yadro/config"
	"github.com/Killazius/yadro/event"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// Proc - структура процессора, который связывает события и участников гонки, отвечает за статистику
type Proc struct {
	cfg         *config.Config                 // Конфигурация гонки
	competitors map[int]*competitor.Competitor // Текущие участники
	logger      io.Writer                      // Интерфейс для логирования
}

// New отвечает за создание нового объекта процессора, если не задан вывод, то берет стандартный вывод в консоль
func New(cfg *config.Config, output io.Writer) *Proc {
	if output == nil {
		output = os.Stdout
	}
	return &Proc{
		cfg:         cfg,
		competitors: make(map[int]*competitor.Competitor),
		logger:      output,
	}
}

// logEvent логирует события связанные с гонкой
func (p *Proc) logEvent(t time.Time, msg string) {
	fmt.Fprintf(p.logger, "[%s] %s\n", t.Format("15:04:05.000"), msg)
}

// logStat логирует статистику участника гонки
func (p *Proc) logStat(msg string) {
	fmt.Fprintf(p.logger, "%s\n", msg)
}

// getCompetitor возвращает участника с указанным ID.
func (p *Proc) getCompetitor(id int) *competitor.Competitor {
	if _, ok := p.competitors[id]; !ok {
		p.competitors[id] = &competitor.Competitor{
			ID:            id,
			LapStartTimes: make([]time.Time, 0),
			LapTimes:      make([]time.Duration, 0),
		}
	}
	return p.competitors[id]
}

// ProcessEvents обрабатывает список событий, проверяет дисквалификации и собирает статистику
func (p *Proc) ProcessEvents(events []*event.Event) {
	for _, e := range events {
		p.processEvent(e)
	}
	p.checkDisqual()
	p.getStats()
}

// processEvent обрабатывает одно событие и обновляет состояние участника
func (p *Proc) processEvent(e *event.Event) {
	c := p.getCompetitor(e.CompetitorID)

	switch e.ID {
	case event.Register:
		c.Registered = true
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) registered", c.ID))

	case event.SetStartTime:
		t, err := time.Parse("15:04:05.000", e.ExtraParams)
		if err != nil {
			panic(err)
		}
		c.PlannedStart = t
		c.StartSet = true
		p.logEvent(e.Time, fmt.Sprintf("The start time for the competitor(%d) was set by a draw to %v",
			c.ID, c.PlannedStart.Format("15:04:05.000")))

	case event.OnStartLine:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) is on the start line", c.ID))

	case event.Started:
		windowStart := c.PlannedStart.Add(-30 * time.Second)
		windowEnd := c.PlannedStart.Add(30 * time.Second)
		if e.Time.Before(windowStart) || e.Time.After(windowEnd) {
			c.Disqualified = true
			p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) is disqualified", c.ID))
			return
		}
		c.Started = true
		c.ActualStart = e.Time
		c.LapStartTimes = append(c.LapStartTimes, e.Time)
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) has started", c.ID))

	case event.OnFiringRange:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) is on the firing range(%s)", c.ID, e.ExtraParams))

	case event.Hit:
		c.Hits++
		p.logEvent(e.Time, fmt.Sprintf("The target(%s) has been hit by competitor(%d)", e.ExtraParams, c.ID))

	case event.LeftFiringRange:
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) left the firing range", c.ID))

	case event.EnteredPenaltyLaps:
		c.LapPenaltyStartTimes = append(c.LapStartTimes, e.Time)
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) entered the penalty laps", c.ID))

	case event.LeftPenaltyLaps:
		if len(c.LapPenaltyStartTimes) > 0 {
			penaltyStart := c.LapPenaltyStartTimes[len(c.LapPenaltyStartTimes)-1]
			c.PenaltyTime += e.Time.Sub(penaltyStart)
			p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) left the penalty laps", c.ID))
		}

	case event.EndMainLap:
		if len(c.LapStartTimes) > 0 {
			lastLapStart := c.LapStartTimes[len(c.LapStartTimes)-1]
			lapTime := e.Time.Sub(lastLapStart)
			c.LapTimes = append(c.LapTimes, lapTime)

			if len(c.LapTimes) == p.cfg.Laps {
				c.Finished = true
				c.EndTime = e.Time
				p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) has finished", c.ID))
			} else {
				p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) ended the main lap", c.ID))
			}
		}

	case event.CantContinue:
		c.Finished = false
		c.EndTime = e.Time
		p.logEvent(e.Time, fmt.Sprintf("The competitor(%d) can`t continue: %s", c.ID, e.ExtraParams))
	default:
		panic("unhandled default case")
	}
}

// checkDisqual проверяет и отмечает участников, которые не стартовали вовремя
func (p *Proc) checkDisqual() {
	for _, c := range p.competitors {
		if c.Registered && !c.Started && time.Now().After(c.PlannedStart.Add(30*time.Second)) {
			c.Disqualified = true
			p.logEvent(time.Now(), fmt.Sprintf("The competitor(%d) is disqualified", c.ID))
		}
	}
}

// getStats сортирует участников и формирует статистику по результатам гонки
func (p *Proc) getStats() {
	sortedCompetitors := make([]*competitor.Competitor, 0, len(p.competitors))
	for _, c := range p.competitors {
		sortedCompetitors = append(sortedCompetitors, c)
	}

	sort.Slice(sortedCompetitors, func(i, j int) bool {
		c1, c2 := sortedCompetitors[i], sortedCompetitors[j]

		if c1.Finished && !c2.Finished {
			return true
		}
		if !c1.Finished && c2.Finished {
			return false
		}

		return c1.EndTime.Sub(c1.ActualStart) < c2.EndTime.Sub(c2.ActualStart)
	})

	for _, c := range sortedCompetitors {
		builder := strings.Builder{}

		switch {
		case c.Disqualified:
			builder.WriteString("[Disqualified] ")
		case !c.Finished:
			builder.WriteString("[NotFinished] ")
		default:
			builder.WriteString("[Finished] ")
		}

		builder.WriteString(fmt.Sprintf("%d ", c.ID))
		builder.WriteString("[")

		for i, lap := range c.LapTimes {
			if i > 0 {
				builder.WriteString(", ")
			}
			speed := float64(p.cfg.LapLen) / lap.Seconds()
			hours := int(lap.Hours())
			minutes := int(lap.Minutes()) % 60
			seconds := int(lap.Seconds()) % 60
			millis := int(lap.Milliseconds()) % 1000
			builder.WriteString(fmt.Sprintf("{%02d:%02d:%02d.%03d, %.3f}",
				hours, minutes, seconds, millis, speed))
		}
		builder.WriteString(strings.Repeat(" {,}", p.cfg.Laps-len(c.LapTimes)))
		builder.WriteString("] ")

		if c.PenaltyTime > 0 {
			hours := int(c.PenaltyTime.Hours())
			minutes := int(c.PenaltyTime.Minutes()) % 60
			seconds := int(c.PenaltyTime.Seconds()) % 60
			millis := int(c.PenaltyTime.Milliseconds()) % 1000

			missedShots := p.cfg.FiringLines*5 - c.Hits
			if missedShots > 0 {
				penaltyDistance := float64(p.cfg.PenaltyLen * missedShots)
				speed := penaltyDistance / c.PenaltyTime.Seconds()
				builder.WriteString(fmt.Sprintf("{%02d:%02d:%02d.%03d, %.3f} ",
					hours, minutes, seconds, millis, speed))
			} else {
				builder.WriteString("{00:00:00.000, 0.000} ")
			}
		} else {
			builder.WriteString("{00:00:00.000, 0.000} ")
		}

		builder.WriteString(fmt.Sprintf("%d/%d", c.Hits, p.cfg.FiringLines*5))

		p.logStat(builder.String())
	}
}
