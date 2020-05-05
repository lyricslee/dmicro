package capx

import (
	"context"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/client"

	"dmicro/common/capx/model"
	"dmicro/common/log"
)

var engine *xorm.Engine
var b broker.Broker

// 本地事务消息表
func TxStorePublished(session *xorm.Session, id int64, topic string, pb interface{}) (err error) {
	var msg []byte
	if msg, err = proto.Marshal(pb.(proto.Message)); err != nil {
		return err
	}

	name := proto.MessageName(pb.(proto.Message))
	v := model.Published{
		Id:     id,
		Topic:  topic,
		Name:   name,
		Msg:    msg,
		Status: 0,
	}

	if _, err = session.InsertOne(&v); err != nil {
		return
	}

	return
}

// 存储从 MQ 中收到的事务消息
func StoreReceived(id int64, topic string, pb interface{}) (err error) {
	if ok, _ := engine.Exist(&model.Received{Id: id}); ok {
		return fmt.Errorf("msg id=%d exist", id)
	}
	var msg []byte
	if msg, err = proto.Marshal(pb.(proto.Message)); err != nil {
		return err
	}

	name := proto.MessageName(pb.(proto.Message))
	v := model.Received{
		Id:     id,
		Topic:  topic,
		Name:   name,
		Msg:    msg,
		Status: 0,
	}

	if _, err = engine.InsertOne(&v); err != nil {
		return
	}

	return
}

func Publish(id int64, topic string, msg interface{}) error {
	log.Debugf("publish topic %s", topic)

	// PubSub at the client/server level works much like RPC but for async comms.
	// 因为 go-micro 是事件驱动实现的异步消息框架，默认 client 是 RPC 调用但是异步的。
	// https://github.com/micro/examples/tree/master/pubsub
	// PubSub 可以指定 topic + 对应的 client (RPC Broker 都可以)
	p := micro.NewEvent(topic, client.DefaultClient) // RPC
	// 投递的时候扔到 MQ
	if err := p.Publish(context.Background(), msg); err != nil {
		log.Error("publish err:", err)
		updatePublished(id, map[string]interface{}{"status": 2})
		return err
	} else {
		log.Debugf("Published %v", msg)
		updatePublished(id, map[string]interface{}{"status": 1})
		return nil
	}
}

// 分布式事务表分了两种：
// 1. sending() 从 Mysql 表中读取要发送的 Topic 数据，然后发送 topic 的 RPC 请求。
// 2. consuming() 从 Mysql 表中读取别的服务写入的 Topic 数据，然后调用本地对应的处理函数来处理。
func Init(e *xorm.Engine) {
	engine = e
	go sending()
	go consuming()
}

func updatePublished(id int64, m map[string]interface{}) error {
	_, err := engine.Table(new(model.Published)).ID(id).Update(m)
	return err
}

func updateReceived(id int64, m map[string]interface{}) error {
	_, err := engine.Table(new(model.Received)).ID(id).Update(m)
	return err
}

func consumed(id int64, status int) {
	updateReceived(id, map[string]interface{}{"status": status})
}

func TxConsumed(session *xorm.Session, id int64) error {
	_, err := session.Table(new(model.Received)).ID(id).Update(map[string]interface{}{"status": 1})
	return err
}

type ConsumerFn func(proto.Message) error

// 这里的 consumers topic 和对应的消费函数是保存在内存中的
var consumers = map[string]ConsumerFn{}

func RegisterConsumer(name string, fn ConsumerFn) {
	consumers[name] = fn
}
