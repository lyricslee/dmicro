package broker

import (
	"dmicro/srv/ums/internal/config"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/nats"
	"sync"
)

var (
	b     broker.Broker
	onceB sync.Once
)

func GetBroker() broker.Broker {
	onceB.Do(func() {
		b = nats.NewBroker(broker.Addrs(config.NatsBroker.Addrs...))
	})

	return b
}