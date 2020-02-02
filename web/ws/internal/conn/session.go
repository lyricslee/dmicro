package conn

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"dmicro/common/constant"
	"dmicro/common/log"
	"dmicro/web/ws/internal/config"
	"dmicro/web/ws/internal/dao"
)

const (
	heartbeatInterval = 60
)

type wsMessage struct {
	typ int
	buf []byte
}

type session struct {
	sync.Mutex
	appid    int    // appid
	platform int    // 终端类型 1: PHONE, 2: PC, 3: WEB 4: MINIAPP
	uid      uint64 // 用户id
	conn     *websocket.Conn
	rchan    chan *wsMessage
	wchan    chan *wsMessage
	echan    chan interface{}
	closed   bool
	ht       time.Time // 最近一次心跳时间
}

func NewSession(appid, platform int, uid uint64, conn *websocket.Conn) (sess *session) {
	sess = &session{
		appid:    appid,
		platform: platform,
		uid:      uid,
		conn:     conn,
		rchan:    make(chan *wsMessage, 10),
		wchan:    make(chan *wsMessage, 10),
		echan:    make(chan interface{}),
		ht:       time.Now(),
	}

	return
}

// 读websocket
func (sess *session) readLoop() {
	defer sess.Close()
	for {
		typ, data, err := sess.conn.ReadMessage()
		if err != nil {
			return
		}

		if !(typ == websocket.TextMessage || typ == websocket.BinaryMessage) {
			continue
		}

		msg := &wsMessage{
			typ: typ,
			buf: data,
		}
		select {
		case sess.rchan <- msg:
		case <-sess.echan:
			return
		}
	}
}

func (sess *session) writeLoop() {
	for {
		select {
		case msg := <-sess.wchan:
			if err := sess.conn.WriteMessage(msg.typ, msg.buf); err != nil {
				sess.Close()
				return
			}
		case <-sess.echan:
			return
		}
	}
}

// 发送消息
func (sess *session) WriteBinaryMessage(buf []byte) (err error) {
	msg := &wsMessage{typ: websocket.BinaryMessage, buf: buf}
	select {
	case sess.wchan <- msg:
		IncSendSucc()
	case <-sess.echan:
		err = ErrConnectionLost
	default:
		err = ErrSendQueueFull
		IncSendFailed()
	}
	return
}

func (sess *session) WriteTextMessage(buf []byte) (err error) {
	msg := &wsMessage{typ: websocket.TextMessage, buf: buf}
	select {
	case sess.wchan <- msg:
		IncSendSucc()
	case <-sess.echan:
		err = ErrConnectionLost
	default:
		err = ErrSendQueueFull
		IncSendFailed()
	}
	return
}

// 读取消息
func (sess *session) ReadMessage() (msg *wsMessage, err error) {
	select {
	case msg = <-sess.rchan:
	case <-sess.echan:
		err = ErrConnectionLost
	}
	return
}

// 关闭连接
func (sess *session) Close() {
	sess.conn.Close()

	sess.Lock()
	defer sess.Unlock()

	if !sess.closed {
		sess.closed = true
		close(sess.echan)
	} else {
		return
	}

	server.Del(sess)
	log.Debug("连接关闭")

	client := dao.GetRedisClient()
	key := fmt.Sprintf(constant.REDIS_KEY_CONNID, sess.appid, sess.uid, sess.platform)
	client.Del(key)

}

func (sess *session) IsAlive() bool {
	sess.Lock()
	defer sess.Unlock()

	now := time.Now()
	if sess.closed || now.Sub(sess.ht) > time.Duration(2*heartbeatInterval)*time.Second {
		return false
	}
	return true
}

// 更新心跳
func (sess *session) KeepAlive() {
	sess.Lock()
	defer sess.Unlock()

	sess.ht = time.Now()
}

func (sess *session) Start() {
	log.Debug("连接建立")

	// TODO: auth
	server.Add(sess)

	client := dao.GetRedisClient()
	key := fmt.Sprintf(constant.REDIS_KEY_CONNID, sess.appid, sess.uid, sess.platform)
	if err := client.Set(key, config.GateId, time.Duration(2*heartbeatInterval)*time.Second).Err(); err != nil {
		log.Error(err)
		sess.Close()
		return
	}

	connId := fmt.Sprintf(constant.REDIS_KEY_CONNID, sess.appid, sess.uid, sess.platform)
	log.Debug(connId)

	go sess.readLoop()
	go sess.writeLoop()
	go sess.healthCheck()

	for {
		msg, err := sess.ReadMessage()
		if err != nil {
			sess.Close()
			return
		}
		sess.handleMessage(msg)
	}
}

func (sess *session) handleMessage(msg *wsMessage) {
	if msg.typ == websocket.TextMessage {
		sess.handleTextMessage(msg)
	} else if msg.typ == websocket.BinaryMessage {
		sess.handleBinaryMessage(msg)
	}
}

func (sess *session) handleTextMessage(msg *wsMessage) {
	// 解析消息体
	req, err := DecodeJSON(msg.buf)
	log.Debug(string(msg.buf))
	if err != nil {
		log.Error(err)
		return
	}
	req.AppId = int32(sess.appid)
	req.Uid = sess.uid
	req.Platform = int32(sess.platform)
	if req.Type == constant.PING {
		sess.handlePing(req)
		return
	}
	if h := MessageHandlers[int(req.Type)]; h == nil {
		log.Error(err)
		return
	} else {
		h(sess, req, 1)
	}
}

func (sess *session) handleBinaryMessage(msg *wsMessage) {
	if len(msg.buf) < 32 {
		log.Error("Invalid packet")
		return
	}
	m, err := DecodeBinary(msg.buf)
	if err != nil {
		log.Error(err)
		return
	}

	if m.Type == constant.PING {
		sess.handlePing(m)
		return
	}
	if h := MessageHandlers[int(m.Type)]; h == nil {
		log.Debug(fmt.Sprintf("Unknown message type=%d", m.Type))
		return
	} else {
		h(sess, m, 2)
	}
}

func (sess *session) healthCheck() {
	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tick.C:
			if !sess.IsAlive() {
				sess.Close()
				return
			}
		case <-sess.echan:
			tick.Stop()
			return
		}
	}
}

type Ping struct {
}

type Pong struct {
}

func (sess *session) handlePing(req *Message) (rsp *Message, err error) {
	sess.KeepAlive()

	b, err := json.Marshal(Pong{})
	if err != nil {
		return
	}

	rsp = &Message{
		Type:    constant.PONG,
		Seq:     req.Seq,
		Payload: json.RawMessage(b),
	}

	buf, err := json.Marshal(*rsp)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err = sess.WriteTextMessage(buf); err != nil {
		if err != ErrSendQueueFull {
			sess.Close()
			return
		}
	}

	client := dao.GetRedisClient()
	key := fmt.Sprintf(constant.REDIS_KEY_CONNID, sess.appid, sess.uid, sess.platform)
	if err = client.Set(key, config.GateId, time.Duration(2*heartbeatInterval)*time.Second).Err(); err != nil {
		log.Error(err)
		return
	}

	return
}
