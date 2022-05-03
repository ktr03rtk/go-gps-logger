package gps

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGps(t *testing.T) {
	t.Parallel()

	date := time.Date(2022, 5, 3, 0, 9, 0, 0, time.Local)

	tests := []struct {
		name           string
		timestamp      time.Time
		mode           int
		lat            float64
		lon            float64
		alt            float64
		speed          float64
		expectedOutput *Gps
		expectedErr    error
	}{
		{
			"normal case1",
			date,
			3,
			30.11,
			130.11,
			30.11,
			30.11,
			&Gps{Timestamp: date, Mode: 3, Lat: 30.11, Lon: 130.11, Alt: 30.11, Speed: 30.11},
			nil,
		},
		{
			"normal case2",
			time.Date(2022, 5, 3, 0, 9, 0, 0, time.Local),
			3,
			-30.11,
			-130.11,
			30.11,
			30.11,
			&Gps{Timestamp: date, Mode: 3, Lat: -30.11, Lon: -130.11, Alt: 30.11, Speed: 30.11},
			nil,
		},
		{
			"error case1 lat exceed maximum limit",
			date,
			3,
			90.1,
			130.11,
			30.11,
			30.11,
			nil,
			errors.New("failed to satisfy GPS Lat Spec"),
		},
		{
			"error case2 lat exceed minimum limit",
			date,
			3,
			-90.1,
			130.11,
			30.11,
			30.11,
			nil,
			errors.New("failed to satisfy GPS Lat Spec"),
		},
		{
			"error case3 lon exceed maximum limit",
			date,
			3,
			30.11,
			180.1,
			30.11,
			30.11,
			nil,
			errors.New("failed to satisfy GPS Lon Spec"),
		},
		{
			"error case4 lon exceed minimum limit",
			date,
			3,
			30.11,
			-180.1,
			30.11,
			30.11,
			nil,
			errors.New("failed to satisfy GPS Lon Spec"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := NewGps(tt.timestamp, tt.mode, tt.lat, tt.lon, tt.alt, tt.speed)
			if err != nil {
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				} else {
					t.Fatalf("error is not expected but received: %v", err)
				}
			} else {
				assert.Exactly(t, tt.expectedErr, nil, "error is expected but received nil")
				assert.Exactly(t, tt.expectedOutput, output)
			}
		})
	}
}
