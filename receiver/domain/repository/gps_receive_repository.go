//go:generate mockgen -source=gps_receive_repository.go -destination=../../mock/mock_gps_receive_repository.go -package=mock
package repository

import "github.com/ktr03rtk/go-gps-logger/receiver/domain/model"

type GpsReceiveRepository interface {
	Receive() (*model.Gps, error)
}
