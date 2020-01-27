package main

import (
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"

	"dmicro/common/capx"
	"dmicro/common/log"
	"dmicro/common/service"
	"dmicro/pkg/micro/go-plugins/broker/stan"
	"dmicro/pkg/micro/go-plugins/wrapper/trace/opentracing"
	"dmicro/pkg/tracer"
	passport "dmicro/srv/passport/api"
	"dmicro/srv/passport/internal/client"
	"dmicro/srv/passport/internal/config"
	"dmicro/srv/passport/internal/dao"
	"dmicro/srv/passport/internal/handler"
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

	// stan
	b := stan.NewBroker(
		broker.Addrs(config.StanBroker.Addrs...),
		stan.ClientID(getClientID()),
		stan.ClusterID(config.StanBroker.ClusterID),
		stan.ConnectRetry(true),
	)
	opts = append(opts, micro.Broker(b))

	svc.Init(opts...)

	// Client
	client.Init(svc)

	// Register Handler
	passport.RegisterPassportHandler(svc.Server(), handler.GetPassportHandler())

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
}
