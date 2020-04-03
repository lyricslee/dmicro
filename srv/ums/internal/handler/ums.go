package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"

	"dmicro/common/constant"
	"dmicro/common/log"
	"dmicro/pkg/util/convert"
	ums "dmicro/srv/ums/api"
	dbroker "dmicro/srv/ums/internal/broker"
	"dmicro/srv/ums/internal/dao"
)

type UmsHandler struct{}

var (
	umsHandler     *UmsHandler
	onceUmsHandler sync.Once
)

// sync.Once 调用 Do() 函数只完成一次初始化
func GetUmsHandler() *UmsHandler {
	onceUmsHandler.Do(func() {
		umsHandler = &UmsHandler{}
	})
	return umsHandler
}

func (this *UmsHandler) Push(ctx context.Context, req *ums.PushRequest, rsp *ums.PushResponse) error {
	return nil
}

// A2L handler，这里分为了 REQ 类型消息的回应和聊天消息
func (this *UmsHandler) A2L(ctx context.Context, req *ums.A2LRequest, rsp *ums.A2LResponse) (err error) {
	log.Debug("Received Ums.A2L request")
	if req.Type == constant.RSP {
		err = this.handleRsp(ctx, req, rsp)
	} else if req.Type == constant.IM {
		err = this.handleIM(ctx, req, rsp)
	}
	return
}

func (this *UmsHandler) G2L(ctx context.Context, req *ums.G2LRequest, rsp *ums.G2LResponse) error {
	log.Debug("Received Ums.G2L request")
	log.Debug(req)
	// 直接转发给应用服务器
	topic := fmt.Sprintf(constant.TOPIC_L2A_PREFIX, req.Appid)

	p := micro.NewEvent(topic, client.DefaultClient) // TODO
	if err := p.Publish(context.Background(), req); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (this *UmsHandler) handleRsp(ctx context.Context, req *ums.A2LRequest, rsp *ums.A2LResponse) error {

	var (
		proto string
		appid string
		uid   string
		plat  string
		cmd   string
		typ   string
		seq   string
		topic string
	)
	rc := dao.GetClient()
	key := fmt.Sprintf(constant.RedisKeyConnid, req.Appid, req.Uid, req.Platform)
	// 1. 从 redis 中读取到该 user 的 conn key
	result, _ := rc.Get(key).Result()
	gateid, _ := convert.ConvertInt(result)

	// broker 的 topic 格式化，这里是 gateid，"topic.L2G:%d"      // %d为gateid
	// TODO
	topic = fmt.Sprintf(constant.TOPIC_L2G_PREFIX, gateid)

	m := &broker.Message{}
	m.Header = make(map[string]string)
	convert.ConvertAssign(&proto, req.Proto)
	convert.ConvertAssign(&appid, req.Appid)
	convert.ConvertAssign(&uid, req.Uid)
	convert.ConvertAssign(&plat, req.Platform)
	convert.ConvertAssign(&typ, req.Type)
	convert.ConvertAssign(&cmd, req.Cmd)
	convert.ConvertAssign(&seq, req.Seq)
	m.Header["proto"] = proto
	m.Header["appid"] = appid
	m.Header["uid"] = uid
	m.Header["plat"] = plat
	m.Header["type"] = typ
	m.Header["cmd"] = cmd
	m.Header["seq"] = seq
	m.Body = req.Payload

	// UMS 服务把收到的 REQ 消息扔到了 Broker
	b := dbroker.GetBroker()
	b.Publish(topic, m)

	return nil
}

func (this *UmsHandler) handleIM(ctx context.Context, req *ums.A2LRequest, rsp *ums.A2LResponse) error {
	return nil
}
