package broker

import (
	"context"
	"dmicro/common/constant"
	"dmicro/common/log"

	topic "dmicro/common/proto/topic"
	"dmicro/srv/chat/internal/client"
	ums "dmicro/srv/ums/api"
)

func HandleL2A(ctx context.Context, msg *topic.L2A) (err error) {
	log.Debug("HandleL2A ...")
	// broker 消息过来有 REQ 和 IM 聊天两种
	if msg.Type == constant.REQ {
		err = handleReq(ctx, msg)
	} else if msg.Type == constant.IM {
		err = handleIM(ctx, msg)
	}

	return
}

// REQ 消息类型，完成具体的业务逻辑之后扔给 UMS 服务。
func handleReq(ctx context.Context, msg *topic.L2A) (err error) {
	req := &ums.A2LRequest{}
	req.Proto = msg.Proto
	req.Appid = msg.Appid
	req.Uid = msg.Uid
	req.Platform = msg.Platform
	req.Type = constant.RSP

	// TODO: 实际应用中要解析Cmd做不同的业务处理，这里简单的回射消息
	req.Cmd = msg.Cmd
	req.Seq = msg.Seq
	req.Payload = msg.Payload
	log.Debug(req)
	_, err = client.UmsClient.A2L(context.Background(), req)
	return
}

func handleIM(ctx context.Context, msg *topic.L2A) (err error) {
	// TODO
	return
}
