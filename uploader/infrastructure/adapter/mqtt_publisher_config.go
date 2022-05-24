package adapter

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type config struct {
	serverURL *url.URL // MQTT server URL
	clientID  string   // Client ID to use when connecting to server
	topic     string   // Topic to subscribe to
	qos       byte     // QOS to use when subscribing

	keepAlive         uint16        // seconds between keepalive packets
	connectRetryDelay time.Duration // Period between connection attempts
}

func getConfig() (config, error) {
	var cfg config
	var err error

	host, ok := os.LookupEnv("MQTT_HOST")
	if !ok {
		return config{}, errors.New("env MQTT_HOST is not found")
	}

	port, ok := os.LookupEnv("MQTT_PORT")
	if !ok {
		return config{}, errors.New("env MQTT_PORT is not found")
	}

	if cfg.serverURL, err = url.Parse("mqtts://" + host + ":" + port); err != nil {
		return config{}, errors.Wrapf(err, "failed to parse mqtt url")
	}

	if cfg.clientID, ok = os.LookupEnv("MQTT_CLIENT"); !ok {
		return config{}, errors.New("env MQTT_CLIENT is not found")
	}

	if cfg.topic, ok = os.LookupEnv("TOPIC"); !ok {
		return config{}, errors.New("env TOPIC is not found")
	}

	sQos, ok := os.LookupEnv("QOS")
	if !ok {
		return config{}, errors.New("env QOS is not found")
	}

	iQos, err := strconv.Atoi(sQos)
	if err != nil {
		return config{}, errors.New("env QOS is not integer")
	}
	cfg.qos = byte(iQos)

	sKa, ok := os.LookupEnv("KEEP_ALIVE")
	if !ok {
		return config{}, errors.New("env KEEP_ALIVE is not found")
	}

	iKa, err := strconv.Atoi(sKa)
	if err != nil {
		return config{}, errors.New("env KEEP_ALIVE is not integer")
	}
	cfg.keepAlive = uint16(iKa)

	sCr, ok := os.LookupEnv("CONNECT_RETRY_DELAY")
	if !ok {
		return config{}, errors.New("env CONNECT_RETRY_DELAY is not found")
	}

	iCr, err := strconv.Atoi(sCr)
	if err != nil {
		return config{}, errors.New("env CONNECT_RETRY_DELAY is not integer")
	}

	cfg.connectRetryDelay = time.Duration(iCr) * time.Millisecond

	return cfg, nil
}

func newTLSConfig() (*tls.Config, error) {
	caPath, ok := os.LookupEnv("CA_PATH")
	if !ok {
		return nil, errors.New("env CA_PATH is not found")
	}

	certPath, ok := os.LookupEnv("CERT_PATH")
	if !ok {
		return nil, errors.New("env CERT_PATH is not found")
	}

	keyPath, ok := os.LookupEnv("KEY_PATH")
	if !ok {
		return nil, errors.New("env KEY_PATH is not found")
	}

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read ca")
	}

	if ok = certpool.AppendCertsFromPEM(pemCerts); !ok {
		return nil, errors.New("failed to append ca")
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.New("failed to load client key pair")
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, errors.New("failed to parse cert")
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}, nil
}
