package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"dmicro/common/log"
	"dmicro/pkg/util/convert"
	"dmicro/web/ws/internal/conn"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// for websocket test
func say(w http.ResponseWriter, r *http.Request) {
	// Upgrade request to websocket
	// 更新 request 参数到 ws 结构体
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// 创建一个 goroutine 来处理 ws 请求
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Info("ReadMessage: ", err)
				return
			}
			log.Info("msg: ", string(msg))
			ws.WriteMessage(websocket.TextMessage, msg)
		}

	}()

}

func join(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Header)

	// 解析 Post 请求 Form 表单
	r.ParseForm()

	var (
		appid int
		uid   uint64
		plat  int
	)
	// 数据类型转换
	convert.ConvertAssign(&appid, r.Form.Get("appid"))
	convert.ConvertAssign(&uid, r.Form.Get("uid"))
	convert.ConvertAssign(&plat, r.Form.Get("plat"))

	if appid == 0 || uid == 0 || plat == 0 {
		log.Errorf("appid uid plat can't be zero appid=%d uid=%d plat=%d", appid, uid, plat)
		return
	}

	log.Debugf(" appid=%d uid=%d plat=%d", appid, uid, plat)
	// Upgrade request to websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Upgrade: ", err)
		return
	}

	// 踢掉旧连的连接
	if sess := conn.GetServer().Get(fmt.Sprintf("%d:%d:%d", appid, uid, plat)); sess != nil {
		sess.Close()
	}
    // 踢掉就的连接后创建新的连接，创建新的 session_id 并且记录在 redis
	sess := conn.NewSession(appid, plat, uid, ws)
	sess.Start()
}
