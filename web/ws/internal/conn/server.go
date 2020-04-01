package conn

import (
	"bytes"
	"dmicro/common/constant"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
)

type Server struct {
	sync.RWMutex
	users    map[string]*session
	sessions map[*session]bool
}

// sync.Once 保证这个 server 是单例的，不会被创建多次。
var (
	server     *Server
	onceServer sync.Once
)

// 创建 onceServer 的时候初始化，sessions 和 users 数据
func GetServer() *Server {
	onceServer.Do(func() {
		server = &Server{
			sessions: make(map[*session]bool),
			users:    make(map[string]*session),
		}
	})

	return server
}

// 通过 id 查找到对应 user 的 sess 然后写消息。
func WriteTextMessage(m *Message) (err error) {
	id := fmt.Sprintf("%d:%d:%d", m.AppId, m.Uid, m.Platform)
	if sess := server.users[id]; sess != nil {
		buf, _ := json.Marshal(m)
		err = sess.WriteTextMessage(buf)
	}
	return
}

// 写二进制消息
func WriteBinaryMessage(m *Message) (err error) {
	var (
		n int32 = 0
	)
	id := fmt.Sprintf("%d:%d:%d", m.AppId, m.Uid, m.Platform)
	if sess := server.users[id]; sess != nil {
		buf := &bytes.Buffer{}
		n = int32(HeaderLen + len(m.Payload))
		err = binary.Write(buf, binary.BigEndian, &n)
		err = binary.Write(buf, binary.BigEndian, &m.AppId)
		err = binary.Write(buf, binary.BigEndian, &m.Uid)
		err = binary.Write(buf, binary.BigEndian, &m.Platform)
		err = binary.Write(buf, binary.BigEndian, &m.Type)
		err = binary.Write(buf, binary.BigEndian, &m.Cmd)
		err = binary.Write(buf, binary.BigEndian, &m.Seq)

		buf.Write(m.Payload)
		err = sess.WriteBinaryMessage(buf.Bytes())
	}
	return
}

func PushTextMessage(appid int32, uid uint64, platform int32, payload string) (err error) {
	id := fmt.Sprintf("%d:%d:%d", appid, uid, platform)
	if sess := server.users[id]; sess != nil {
		m := Message{}
		m.Type = constant.PUSH
		m.Payload = []byte(payload)
		buf, _ := json.Marshal(m)
		err = sess.WriteTextMessage(buf)
	}
	return
}

func PushBinaryMessage(appid int32, uid uint64, platform int32, payload []byte) (err error) {
	var (
		typ int32 = constant.PUSH
		cmd int32 = 0
		seq int32 = 0
	)
	id := fmt.Sprintf("%d:%d:%d", appid, uid, platform)
	if sess := server.users[id]; sess != nil {
		buf := &bytes.Buffer{}
		n := HeaderLen + len(payload)
		err = binary.Write(buf, binary.BigEndian, &n)
		err = binary.Write(buf, binary.BigEndian, &appid)
		err = binary.Write(buf, binary.BigEndian, &uid)
		err = binary.Write(buf, binary.BigEndian, &platform)
		err = binary.Write(buf, binary.BigEndian, &typ)
		err = binary.Write(buf, binary.BigEndian, &cmd)
		err = binary.Write(buf, binary.BigEndian, &seq)
		buf.Write(payload)
		err = sess.WriteBinaryMessage(buf.Bytes())
	}
	return
}

func Response(msg *Message) {
	id := fmt.Sprintf("%d:%d:%d", msg.AppId, msg.Uid, msg.Platform)
	if sess := server.users[id]; sess != nil {
		sess.WriteTextMessage(msg.Payload)
	}
}

// 建立连接的时候把用户 sess 信息添加到 users 列表中
func (s *Server) Add(sess *session) {
	s.Lock()
	defer s.Unlock()
	id := fmt.Sprintf("%d:%d:%d", sess.appid, sess.uid, sess.platform)
	s.users[id] = sess
	IncOnlineNum()
}

// 删除离线用户
func (s *Server) Del(sess *session) {
	s.Lock()
	defer s.Unlock()
	delete(s.sessions, sess)
	id := fmt.Sprintf("%d:%d:%d", sess.appid, sess.uid, sess.platform)
	delete(s.users, id)
	DecOnlineNum()
}

// 通过 id 或者 session
func (s *Server) Get(id string) (sess *session) {
	s.Lock()
	defer s.Unlock()
	sess = s.users[id]
	return
}

// 统计在在线人数
type Stats struct {
	OnlineNum     int64 `json:"online_num"`
	SendSuccNum   int64 `json:"send_succ_num"`
	SendFailedNum int64 `json:"send_failed_num"`
}

var (
	stats Stats
)

func IncOnlineNum() {
	atomic.AddInt64(&stats.OnlineNum, 1)
}

func DecOnlineNum() {
	atomic.AddInt64(&stats.OnlineNum, -1)
}

func IncSendSucc() {
	atomic.AddInt64(&stats.SendSuccNum, 1)
}

func DecSendSucc() {
	atomic.AddInt64(&stats.SendSuccNum, -1)
}

func IncSendFailed() {
	atomic.AddInt64(&stats.SendFailedNum, 1)
}

func DecSendFailed() {
	atomic.AddInt64(&stats.SendFailedNum, -1)
}

func (stats *Stats) Dump() (data []byte, err error) {
	return json.Marshal(stats)
}
