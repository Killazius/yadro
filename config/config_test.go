package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMustLoad(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		cfgJSON := `{
		"laps": 3,
		"lapLen": 1000,
		"penaltyLen": 150,
		"firingLines": 2,
		"start": "10:00",
		"startDelta": "00:30"
	}`
		temp, err := os.CreateTemp("", "cfg_*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(temp.Name())
		if _, err = temp.WriteString(cfgJSON); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		if err = temp.Close(); err != nil {
			t.Fatalf("Failed to close temp file: %v", err)
		}
		cfg := MustLoad(temp.Name())
		require.NotNil(t, cfg)
		require.Equal(t, 3, cfg.Laps)
		require.Equal(t, 1000, cfg.LapLen)
		require.Equal(t, 150, cfg.PenaltyLen)
		require.Equal(t, 2, cfg.FiringLines)
		require.Equal(t, "10:00", cfg.Start)
		require.Equal(t, "00:30", cfg.StartDelta)
	})

	errorCases := []struct {
		name     string
		config   string
		expected string
	}{
		{
			name:     "file not exists",
			config:   "",
			expected: "config file does not exist",
		},
		{
			name:     "invalid JSON",
			config:   `{"laps": 3,}`,
			expected: "error loading config",
		},
		{
			name:     "invalid field type",
			config:   `{"laps": "three"}`,
			expected: "error loading config",
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			var temp *os.File
			var err error

			if tc.config != "" {
				temp, err = os.CreateTemp("", "invalid_*.json")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(temp.Name())

				if _, err = temp.WriteString(tc.config); err != nil {
					t.Fatal(err)
				}
				temp.Close()
			}

			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic but got none")
				}
			}()
			if tc.config == "" {
				MustLoad("file.json")
			} else {
				MustLoad(temp.Name())
			}
		})
	}
}
