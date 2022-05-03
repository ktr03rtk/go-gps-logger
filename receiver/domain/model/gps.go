package gps

import (
	"time"

	"github.com/pkg/errors"
)

type Gps struct {
	Timestamp time.Time `json:"timestamp"`
	Mode      int       `json:"mode"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Alt       float64   `json:"alt"`
	Speed     float64   `json:"speed"`
}

const (
	MINIMUM_LAT = -90
	MAXIMUM_LAT = 90
	MINIMUM_LON = -180
	MAXIMUM_LON = 180
)

func NewGps(time time.Time, mode int, lat, lon, alt, speed float64) (*Gps, error) {
	g := &Gps{
		Timestamp: time,
		Mode:      mode,
		Lat:       lat,
		Lon:       lon,
		Alt:       alt,
		Speed:     speed,
	}

	if err := gpsSpecSatisfied(*g); err != nil {
		return nil, err
	}

	return g, nil
}

func gpsSpecSatisfied(g Gps) error {
	if g.Lat < MINIMUM_LAT || g.Lat > MAXIMUM_LAT {
		return errors.Errorf("failed to satisfy GPS Lat Spec")
	}

	if g.Lon < MINIMUM_LON || g.Lon > MAXIMUM_LON {
		return errors.Errorf("failed to satisfy GPS Lon Spec")
	}

	return nil
}
