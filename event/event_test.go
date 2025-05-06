package event

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		event   *Event
		wantErr bool
	}{
		{
			name: "Register event",
			line: "[09:05:59.867] 1 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 5, 59, 867000000, time.UTC),
				ID:           Register,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "SetStartTime with extra param",
			line: "[09:15:00.841] 2 1 09:30:00.000",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 15, 0, 841000000, time.UTC),
				ID:           SetStartTime,
				CompetitorID: 1,
				ExtraParams:  "09:30:00.000",
			},
			wantErr: false,
		},
		{
			name: "OnStartLine event",
			line: "[09:29:45.734] 3 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 29, 45, 734000000, time.UTC),
				ID:           OnStartLine,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "Started event",
			line: "[09:30:01.005] 4 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 30, 1, 5000000, time.UTC),
				ID:           Started,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "OnFiringRange with target number",
			line: "[09:49:31.659] 5 1 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 49, 31, 659000000, time.UTC),
				ID:           OnFiringRange,
				CompetitorID: 1,
				ExtraParams:  "1",
			},
			wantErr: false,
		},
		{
			name: "Hit target",
			line: "[09:49:33.123] 6 1 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 49, 33, 123000000, time.UTC),
				ID:           Hit,
				CompetitorID: 1,
				ExtraParams:  "1",
			},
			wantErr: false,
		},
		{
			name: "LeftFiringRange event",
			line: "[09:49:38.339] 7 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 49, 38, 339000000, time.UTC),
				ID:           LeftFiringRange,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "EnteredPenaltyLaps event",
			line: "[09:49:55.915] 8 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 49, 55, 915000000, time.UTC),
				ID:           EnteredPenaltyLaps,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "LeftPenaltyLaps event",
			line: "[09:51:48.391] 9 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 51, 48, 391000000, time.UTC),
				ID:           LeftPenaltyLaps,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "EndMainLap event",
			line: "[09:59:03.872] 10 1",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 59, 3, 872000000, time.UTC),
				ID:           EndMainLap,
				CompetitorID: 1,
				ExtraParams:  "",
			},
			wantErr: false,
		},
		{
			name: "CantContinue with reason",
			line: "[09:59:03.872] 11 1 Lost in the forest",
			event: &Event{
				Time:         time.Date(0, 1, 1, 9, 59, 3, 872000000, time.UTC),
				ID:           CantContinue,
				CompetitorID: 1,
				ExtraParams:  "Lost in the forest",
			},
			wantErr: false,
		},
		{
			name:    "invalid line format",
			line:    "[09:05:59.867]",
			event:   nil,
			wantErr: true,
		},
		{
			name:    "invalid time format",
			line:    "[25:05:59.867] 1 1",
			event:   nil,
			wantErr: true,
		},
		{
			name:    "invalid event ID",
			line:    "[09:05:59.867] abc 1",
			event:   nil,
			wantErr: true,
		},
		{
			name:    "invalid competitor ID",
			line:    "[09:05:59.867] 1 xyz",
			event:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.event) {
				t.Errorf("= %v, want %v", got, tt.event)
			}
		})
	}
}
