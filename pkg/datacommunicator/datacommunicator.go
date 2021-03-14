package datacommunicator

import (
	"log"

	gpsd "github.com/koppacetic/go-gpsd"
)

type Communicator struct {
	session       *gpsd.Session
	received_data chan gpsd.TPVReport
}

func NewCommunicator() *Communicator {
	c := new(Communicator)
	session, err := gpsd.Dial(gpsd.DefaultAddress)
	if err != nil {
		log.Fatalf("Failed to connect to GPSD: %s", err)
	}
	c.session = session
	c.received_data = make(chan gpsd.TPVReport)
	return c
}

func (c *Communicator) Communicate() {

	c.session.Subscribe("TPV", func(r interface{}) {
		tpv := r.(*gpsd.TPVReport)
		c.received_data <- *tpv
	})

	c.session.Run()
}

func (c *Communicator) Receive() gpsd.TPVReport {
	data := <-c.received_data
	return data
}

func (c *Communicator) Close() error {
	err := c.session.Close()
	return err
}
