package messaging

import (
	"time"

	"github.com/desulaidovich/balance/internal/utils"
	"github.com/nats-io/nats.go"
)

type NatsConnection struct {
	*nats.Conn
}

func (n *NatsConnection) SendJSON(subject string, data utils.JSONMessage) error {
	encodedCon, err := nats.NewEncodedConn(n.Conn, nats.JSON_ENCODER)

	if err != nil {
		return err
	}

	if err = encodedCon.Publish(subject, &data); err != nil {
		return err
	}

	return nil
}

func NewNatsConnection(conn *nats.Conn) *NatsConnection {
	return &NatsConnection{
		conn,
	}
}

func Connect() (*nats.Conn, error) {
	conn, err := nats.Connect(nats.DefaultURL, nats.Name("Balance"), nats.Timeout(10*time.Second))

	if err != nil {
		return nil, err
	}

	return conn, nil
}
