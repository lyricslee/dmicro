package handler

import (
	"context"
	"sync"

	"github.com/go-xorm/xorm"
	"github.com/golang/protobuf/proto"

	"dmicro/common/capx"
	"dmicro/common/log"
	topic "dmicro/common/proto/topic"
	"dmicro/pkg/tx"
	"dmicro/srv/user/internal/dao"
	"dmicro/srv/user/internal/model"
	"dmicro/srv/user/internal/service"
)

type Subscriber struct {
	engine *xorm.Engine
	svc    *service.UserService
}

var (
	subscriber     *Subscriber
	onceSubscriber sync.Once
)

func GetSubscriber() *Subscriber {
	onceSubscriber.Do(func() {
		subscriber = &Subscriber{
			engine: dao.GetEngine(),
			svc:    service.GetUserService(),
		}
	})
	return subscriber
}

//func NewSubscriber(engine *xorm.Engine, svc *service.UserService) *Subscriber {
//	return &Subscriber{engine: engine, svc: svc}
//}

func (sub *Subscriber) UserCreated() func(ctx context.Context, msg *topic.UserCreated) (err error) {
	return func(ctx context.Context, msg *topic.UserCreated) (err error) {
		log.Debugf("received msg topic=%s", msg.Topic)
		// 事务消息表中存 topic 记录
		if err = capx.StoreReceived(msg.Id, msg.Topic, msg); err != nil {
			// 重复的消息，保存失败，从而保证幂等性
			return
		}

		session := sub.engine.NewSession()
		defer func() {
			session.Close()
		}()
		// 事务：1. 用户表插入数据
		err = tx.ExecWithTransaction(session, func(session *xorm.Session) (err error) {
			if _, err = session.InsertOne(&model.User{Id: msg.Info.Uid, Mobile: msg.Info.Mobile}); err != nil {
				return
			}

			// 2. 设置事务 topic 执行成功了
			if err = capx.TxConsumed(session, msg.Id); err != nil {
				return
			}
			return
		})

		return
	}
}

// 由定时器触发的重新消费函数，来保证事务执行。
func (sub *Subscriber) CapxUserCreated() capx.ConsumerFn {
	return func(pb proto.Message) (err error) {
		log.Debug("CapxUserCreated")
		msg := pb.(*topic.UserCreated)
		return sub.svc.Create(msg.Info.Uid, msg.Info.Mobile)
	}
}
