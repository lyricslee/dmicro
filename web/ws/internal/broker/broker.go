package broker

import (
	"fmt"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/nats"

	"dmicro/common/constant"
	"dmicro/common/log"
	"dmicro/pkg/util/convert"
	"dmicro/web/ws/internal/config"
	"dmicro/web/ws/internal/conn"
)

var b broker.Broker

func Init() {
	b = nats.NewBroker(broker.Addrs(config.NatsBroker.Addrs...))

	if err := b.Init(); err != nil {
		log.Fatalf("Broker init error: %v", err)
	}

	if err := b.Connect(); err != nil {
		log.Fatalf("Broker connect error: %v", err)
	}

	// TODO: config gateid
	topic := fmt.Sprintf(constant.TOPIC_L2G_PREFIX, "1")
	if _, err := b.Subscribe(topic, handleL2G); err != nil {
		log.Fatalf("Broker subscribe error: %v", err)
	}
}

func GetBroker() broker.Broker {
	return b
}

func handleL2G(event broker.Event) (err error) {
	var (
		proto int32
		appid int32
		uid   uint64
		plat  int32
		typ   int32
		cmd   int32
		seq   int32
	)
	convert.ConvertAssign(&proto, event.Message().Header["proto"])
	convert.ConvertAssign(&appid, event.Message().Header["appid"])
	convert.ConvertAssign(&uid, event.Message().Header["uid"])
	convert.ConvertAssign(&plat, event.Message().Header["plat"])
	convert.ConvertAssign(&typ, event.Message().Header["type"])
	convert.ConvertAssign(&cmd, event.Message().Header["cmd"])
	convert.ConvertAssign(&seq, event.Message().Header["seq"])
	m := &conn.Message{}
	m.AppId = appid
	m.Uid = uid
	m.Platform = plat
	m.Type = typ
	m.Cmd = cmd
	m.Seq = seq
	m.Payload = event.Message().Body

	if proto == 1 {
		err = conn.WriteTextMessage(m)
	} else {
		err = conn.WriteBinaryMessage(m)
	}
	return
}
