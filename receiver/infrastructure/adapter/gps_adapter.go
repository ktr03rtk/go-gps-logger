package adapter

import (
	"context"
	"sync"

	gpsd "github.com/koppacetic/go-gpsd"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/model"
	"github.com/ktr03rtk/go-gps-logger/receiver/domain/repository"
	"github.com/pkg/errors"
)

type gpsAdapter struct {
	mu         sync.Mutex
	session    *gpsd.Session
	latestData *model.Gps
	err        error
}

type receivedData struct {
	tpvReport *gpsd.TPVReport
	err       error
}

func NewGpsAdapter(ctx context.Context) (repository.GpsReceiveRepository, error) {
	session, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to gpad")
	}

	c := &gpsAdapter{
		session: session,
	}

	c.connect(ctx)

	return c, nil
}

func (c *gpsAdapter) Receive() (*model.Gps, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.latestData == nil {
		return nil, errors.New("no gps data stored")
	}

	if c.err != nil {
		return nil, c.err
	}

	result := *c.latestData
	c.latestData = nil

	return &result, nil
}

func (c *gpsAdapter) connect(ctx context.Context) {
	dataCh := make(chan receivedData)

	c.session.Subscribe("TPV", func(r interface{}) {
		tpv, ok := r.(*gpsd.TPVReport)
		if !ok {
			dataCh <- receivedData{
				tpvReport: nil,
				err:       errors.New("failed to assert TPV report type"),
			}

			return
		}

		dataCh <- receivedData{
			tpvReport: tpv,
			err:       nil,
		}
	})

	c.session.Run()

	go func() {
		c.receive(ctx, dataCh)
	}()

	return
}

func (c *gpsAdapter) receive(ctx context.Context, dataCh <-chan receivedData) {
	for {
		select {
		case <-ctx.Done():
			c.session.Close()

			return
		case r := <-dataCh:
			c.mu.Lock()
			c.latestData = convert(*r.tpvReport)
			c.err = r.err
			c.mu.Unlock()
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
