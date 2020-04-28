package main

import (
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"

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

// 生成 uuid 之后保存在 cid 文件中，这样每次客户端启动都是相同的 cid，broker 就可以找到相同的 clientId 了。
func getClientID() (cid string) {
	var (
		f   *os.File
		err error
		b   []byte
	)

	// cid 是个文件，内容：f509581d-34fa-4123-91a3-75e12a89f843
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

	// passport main.go 与之前 wed/dd 的基本一样，就是初始化配置 服务等。
	svc := service.NewService()

	var opts []micro.Option
	// Tracer init
	t, err := tracer.Init(config.Micro.ServerName, config.Tracer.Addr)
	if err != nil {
		log.Error(err)
	}
	// trace Inject 的地方
	opts = append(opts, micro.WrapHandler(opentracing.NewHandlerWrapper(t)))

	// stan
	// 这里多了一个 broker 的初始化，连接到 broker (nats) 集群上。
	b := stan.NewBroker(
		broker.Addrs(config.StanBroker.Addrs...),
		stan.ClientID(getClientID()), // 本客户端的唯一 id
		stan.ClusterID(config.StanBroker.ClusterID),
		stan.ConnectRetry(true),
	)
	opts = append(opts, micro.Broker(b))

	svc.Init(opts...)

	// Client
	client.Init(svc)

	// Register Handler
	// 初始化 passport 的 handler，也就是 RPC 调用方法对应的处理函数。
	passport.RegisterPassportHandler(svc.Server(), handler.GetPassportHandler())

	// ORM Mysql 初始化
	// capx 模块内部一个 	go sending(） 和一个 go consuming() routine 负载数据库具体读写操作
	capx.Init(dao.GetEngine())

	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
