package conn

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"

	"dmicro/common/constant"
	"dmicro/common/log"
	ums "dmicro/srv/ums/api"
	"dmicro/web/ws/internal/client"
)

const (
	HeaderLen = 32
)

type messageHandler func(*session, *Message, int)

var (
	MessageHandlers = make(map[int]messageHandler)
)

func init() {
	MessageHandlers[constant.REQ] = handleReq
	MessageHandlers[constant.IM] = handleIM
}

type Message struct {
	AppId    int32           `json:"appid,omitempty"`    // appid
	Uid      uint64          `json:"uid,omitempty"`      // 用户ID
	Platform int32           `json:"platform,omitempty"` // 终端类型 1: PHONE, 2: PC, 3: WEB 4: MINIAPP
	Type     int32           `json:"type"`               // type消息类型 1: PING, 2: PONG, 3: REQ, 4: RSP, 5: PUSH, 6: IM
	Cmd      int32           `json:"cmd,omitempty"`      // 命令类型
	Seq      int32           `json:"seq,omitempty"`      // 序号
	Payload  json.RawMessage `json:"payload,omitempty"`  // 数据字段
}

func DecodeJSON(buf []byte) (msg *Message, err error) {
	msg = &Message{}
	if err = json.Unmarshal(buf, msg); err != nil {
		return nil, err
	}

	return
}

func DecodeBinary(b []byte) (m *Message, err error) {

	buf := &bytes.Buffer{}
	buf.Write(b)
	m = &Message{}
	var n int32
	if err := binary.Read(buf, binary.BigEndian, &n); err != nil {
		log.Error(err)
		return nil, err
	}
	if len(b) != int(n) {
		log.Error("Invalid packet")
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &m.AppId); err != nil {
		log.Error(err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &m.Uid); err != nil {
		log.Error(err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &m.Platform); err != nil {
		log.Error(err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &m.Type); err != nil {
		log.Error(err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &m.Cmd); err != nil {
		log.Error(err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &m.Seq); err != nil {
		log.Error(err)
		return nil, err
	}
	m.Payload = make([]byte, len(b)-HeaderLen)
	buf.Read(m.Payload)

	return
}

func handleReq(sess *session, req *Message, mtyp int) {
	r := &ums.G2LRequest{}
	r.Proto = int32(mtyp)
	r.Appid = req.AppId
	r.Uid = req.Uid
	r.Platform = req.Platform
	r.Type = constant.REQ
	r.Cmd = req.Cmd
	r.Seq = req.Seq
	r.Payload = req.Payload

	log.Debug(r)
	client.UmsClient.G2L(context.Background(), r)
}

func handleIM(sess *session, req *Message, mtyp int) {
	// 单聊
	// 群聊

	return
}
