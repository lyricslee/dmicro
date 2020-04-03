package main

import (
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/server"

	"dmicro/common/capx"
	"dmicro/common/constant"
	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/broker/stan"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	user "dmicro/srv/user/api"
	"dmicro/srv/user/internal/config"
	"dmicro/srv/user/internal/dao"
	"dmicro/srv/user/internal/handler"
)

func getClientID() (cid string) {
	var (
		f   *os.File
		err error
		b   []byte
	)

	if f, err = os.OpenFile("cid", os.O_CREATE|os.O_RDWR, 0666); err != nil {
		panic(err)
	}
	defer f.Close()
	if b, err = ioutil.ReadAll(f); err != nil {
		panic(err)
	}
	if len(b) == 0 {
		cid = uuid.New().String()
		if _, err = f.Write([]byte(cid)); err != nil {
			panic(err)
		}
	} else {
		cid = string(b)
	}
	f.Close()
	return
}

func main() {
	// Config
	config.Init()
	// Logger
	log.Init(config.Logger)

	svc := service.NewService()

	var opts []micro.Option
	// Tracer
	t, err := tracer.Init(config.Micro.ServerName, config.Tracer.Addr)
	if err != nil {
		log.Error(err)
	}
	opts = append(opts, micro.WrapHandler(opentracing.NewHandlerWrapper(t)))

	// 1. 初始化 nats broker 连接
	// 设置 nats broker ClientID，有效时长。
	// clientID 和 durableName 对于NATS Streaming非常重要. 要让subscriber重启后能继续
	// 收到重启期间发过来的消息且不重复消息
	b := stan.NewBroker(
		broker.Addrs(config.StanBroker.Addrs...),
		stan.ClientID(getClientID()),
		stan.DurableName(config.StanBroker.DurableName),
		stan.ClusterID(config.StanBroker.ClusterID),
		stan.ConnectRetry(true),
	)
	opts = append(opts, micro.Broker(b))

	svc.Init(opts...)

	// Register Handler
	user.RegisterUserHandler(svc.Server(), handler.GetUserHandler())

	// 2. 初始化完成后订阅消息, 用户创建消息和对应的 handler。
	micro.RegisterSubscriber(
		constant.TOPIC_USER_CREATED,
		svc.Server(),
		handler.GetSubscriber().UserCreated(),
		server.SubscriberQueue(config.StanBroker.Queue), // TODO
	)
	// 注册消费者 用户创建的消息
	capx.RegisterConsumer(constant.TOPIC_USER_CREATED, handler.GetSubscriber().CapxUserCreated())
    // 初始化 MySQL Engine
	capx.Init(dao.GetEngine())

	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}

	//if err := b.Init(); err != nil {
	//	log.Fatalf("Broker init error: %v", err)
	//}
	//
	//if err := b.Connect(); err != nil {
	//	log.Fatalf("Broker connect error: %v", err)
	//}
	//
	//if _, err := b.Subscribe(constant.TOPIC_USER_CREATED, func(event broker.Event) error {
	//	// TODO: reflect
	//	return nil
	//}); err != nil {
	//	log.Fatalf("Broker subscribe error: %v", err)
	//}
}
