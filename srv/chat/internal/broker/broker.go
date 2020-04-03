package broker

import (
	"dmicro/srv/chat/internal/config"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"sync"
)

var (
	b     broker.Broker
	onceB sync.Once
)

// Broker Init
func GetBroker() broker.Broker {
	onceB.Do(func() {
		b = nats.NewBroker(broker.Addrs(config.NatsBroker.Addrs...))
	})

	return b
}
