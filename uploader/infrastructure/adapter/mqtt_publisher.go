package adapter

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/model"
	"github.com/ktr03rtk/go-gps-logger/uploader/domain/repository"
	"github.com/pkg/errors"
)

const (
	connectionWaitTime = 1000 // milliseconds
	timeFormat         = "2006-01-02-15-04"
)

type mqttAdapter struct {
	connectionManager *autopaho.ConnectionManager
	clientID          string
	topic             string
	qos               byte
}

func NewMqttAdapter(ctx context.Context) (repository.PayloadUploadRepository, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	tlsCfg, err := newTLSConfig()
	if err != nil {
		return nil, err
	}

	mqttCfg := getMqttConfig(cfg, tlsCfg)

	cm, err := autopaho.NewConnection(ctx, mqttCfg)
	if err != nil {
		return nil, err
	}

	time.Sleep(connectionWaitTime * time.Millisecond)

	return &mqttAdapter{
		connectionManager: cm,
		clientID:          cfg.clientID,
		topic:             cfg.topic,
		qos:               cfg.qos,
	}, nil
}

func (a *mqttAdapter) Upload(ctx context.Context, payload *model.Payload) ([]model.BaseFilePath, error) {
	if len(payload.FilePaths) == 0 {
		return []model.BaseFilePath{}, nil
	}

	topic, err := a.createTopic(payload)
	if err != nil {
		return nil, err
	}

	if _, err := a.connectionManager.Publish(ctx, &paho.Publish{
		QoS:     a.qos,
		Topic:   topic,
		Payload: payload.Message,
	}); err != nil {
		return nil, errors.Wrapf(err, "failed to publish")
	}

	return payload.FilePaths, nil
}

func getMqttConfig(cfg config, tlsCfg *tls.Config) autopaho.ClientConfig {
	return autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{cfg.serverURL},
		KeepAlive:         cfg.keepAlive,
		ConnectRetryDelay: cfg.connectRetryDelay,
		OnConnectionUp:    func(*autopaho.ConnectionManager, *paho.Connack) { log.Print("mqtt connection up") },
		OnConnectError:    func(err error) { log.Printf("error whilst attempting connection: %s\n", err) },
		Debug:             paho.NOOPLogger{},
		TlsCfg:            tlsCfg,
		ClientConfig: paho.ClientConfig{
			ClientID:      cfg.clientID,
			OnClientError: func(err error) { log.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
}

func (a *mqttAdapter) createTopic(payload *model.Payload) (string, error) {
	t, err := time.Parse(timeFormat, strings.TrimSuffix(string(payload.FilePaths[0]), ".dat"))
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse date")
	}

	topic := a.topic + "/thing=" + a.clientID + t.Format("/year=2006/month=01/day=02/") + string(payload.FilePaths[0])

	return topic, nil
}
