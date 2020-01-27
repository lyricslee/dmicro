package capx

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"

	"dmicro/common/capx/model"
	"dmicro/common/log"
)

func sending() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic recovered: %v", r)
		}
	}()
	tick := time.NewTicker(time.Second * 3)
	for _ = range tick.C {
		m := new(model.Published)

		now := time.Now()

		// 第一次重试间隔5分钟
		rows, err := engine.Where("retries<=100").
			And("status=0 or status=2").
			And("created_at < ?", now.Add(-time.Second*5).Format("2006-01-02 15:04:05")).
			Rows(m)
		if err != nil {
			log.Error(err)
			continue
		}

		for rows.Next() {
			if err := rows.Scan(m); err != nil {
				log.Error(err)
				continue
			}

			// 重试间隔1分钟
			if now.Sub(m.UpdatedAt) < time.Minute {
				continue
			}

			msgType := proto.MessageType(m.Name)
			msg := reflect.Indirect(reflect.New(msgType.Elem())).Addr().Interface().(proto.Message)
			if err := proto.Unmarshal(m.Msg, msg); err != nil {
				log.Error(err)
				continue
			}
			// TODO: 重新投递,retries次数要增加
			if err := Publish(m.Id, m.Topic, msg); err != nil {
				log.Error(err)
				continue
			}
		}

		rows.Close()
	}

}

func consuming() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic recovered: %v", r)
		}
	}()
	tick := time.NewTicker(time.Second * 3)
	for _ = range tick.C {
		m := new(model.Received)

		now := time.Now()

		// 第一次重试间隔5分钟
		rows, err := engine.Where("retries<=100").
			And("status=0 or status=2").
			And("created_at < ?", now.Add(-time.Minute*5)).
			Rows(m)

		if err != nil {
			log.Error(err)
			continue
		}

		for rows.Next() {
			if err := rows.Scan(m); err != nil {
				log.Error(err)
				continue
			}

			// 重试间隔1分钟
			if now.Sub(m.UpdatedAt) < time.Minute {
				continue
			}

			msgType := proto.MessageType(m.Name)
			msg := reflect.Indirect(reflect.New(msgType.Elem())).Addr().Interface().(proto.Message)
			if err := proto.Unmarshal(m.Msg, msg); err != nil {
				log.Error(err)
				continue
			}

			if fn := consumers[m.Topic]; fn != nil {
				if err := fn(msg); err != nil {
					consumed(m.Id, 2)
				} else {
					consumed(m.Id, 1)
				}
			}
		}

		rows.Close()
	}
}
