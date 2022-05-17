package adapter

import (
	"context"
	"fmt"

	gpsd "github.com/koppacetic/go-gpsd"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/pkg/errors"
)

type GpsAdapter struct {
	session      *gpsd.Session
	receivedData chan gpsd.TPVReport
	latestData   *model.Gps
}

func NewGpsAdapter(ctx context.Context) (*GpsAdapter, error) {
	session, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to gpad")
	}

	receivedData := make(chan gpsd.TPVReport)

	c := &GpsAdapter{
		session:      session,
		receivedData: receivedData,
	}

	c.connect(ctx)

	return c, nil
}

func (c *GpsAdapter) connect(ctx context.Context) {
	c.session.Subscribe("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		c.receivedData <- *tpv
	})

	c.session.Run()

	go func() {
		c.receive(ctx)
	}()

	return
}

func (c *GpsAdapter) receive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.session.Close()

			return
		case r := <-c.receivedData:
			c.latestData = convert(r)
		}
	}
}

func convert(r gpsd.TPVReport) *model.Gps {
	data := &model.Gps{
		Timestamp: r.Time,
		Mode:      int(r.Mode),
	}

	if r.Mode == 1 {
		return data
	}

	data.Lat = r.Lat
	data.Lon = r.Lon

	if r.Mode == 2 {
		return data
	}

	data.Alt = r.Alt
	data.Speed = r.Speed

	return data
}

func (c *GpsAdapter) Receive() (*model.Gps, error) {
	if c.latestData == nil {
		return nil, errors.New("no gps data stored")
	}

	result := *c.latestData
	c.latestData = nil
	fmt.Println(result)

	return &result, nil
}
