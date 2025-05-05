package competitor

import "time"

type Competitor struct {
	ID                   int
	Registered           bool
	Disqualified         bool
	PlannedStart         time.Time
	ActualStart          time.Time
	EndTime              time.Time
	StartSet             bool
	Started              bool
	Finished             bool
	Hits                 int
	LapStartTimes        []time.Time
	LapPenaltyStartTimes []time.Time
	LapTimes             []time.Duration
	PenaltyTime          time.Duration
}
