package mq

import (
	"errors"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	ErrNotImplemented = errors.New("Function not implemented")
)

// 订阅得到回应的回调
type CbSubscribe func(subj string, data []byte)

// queue 订阅模式的回调
type CbQueueSubscribe func(subj string, q string, data []byte)

type Eventbus struct {
	addrs []string // 当前连接的地址参数
	conn  *nats.Conn
}

func NewEventbus(addrs []string, username, password string) (*Eventbus, error) {
	nc, err := nats.Connect(strings.Join(addrs, ","),
		nats.NoEcho(),
		nats.Name("Kg Eventbus"),
		nats.Timeout(3*time.Second),
		nats.MaxReconnects(100),
		nats.ReconnectWait(5*time.Second),
		nats.UserInfo(username, password),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			//log.Info("Disconnected to EventServer", zap.String("url", nc.ConnectedUrl()), zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			//log.Info("Reconnect to EventServer", zap.String("url", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			//log.Info("Close connection to Eventbus", zap.String("url", nc.ConnectedUrl()))
		}),
	)
	if nil != err {
		return nil, err
	}
	return &Eventbus{addrs, nc}, nil
}

func (self *Eventbus) GetOriginConn() *nats.Conn {
	return self.conn
}

func (self *Eventbus) AsyncPub(subj string, content []byte, timeout time.Duration) error {
	err := self.conn.Publish(subj, content)
	if err == nil {
		self.conn.FlushTimeout(timeout)
	}
	return err
}

func (self *Eventbus) Sub(subj string, cb CbSubscribe) error {
	_, err := self.conn.Subscribe(subj, func(msg *nats.Msg) {
		cb(msg.Subject, msg.Data)
	})
	if err == nil {
		self.conn.Flush()
	}
	return err
}

func (self *Eventbus) QSub(subj, q string, cb CbQueueSubscribe) error {
	self.conn.QueueSubscribe(subj, q, func(msg *nats.Msg) {
		cb(msg.Subject, msg.Sub.Queue, msg.Data)
	})
	return nil
}

func (self *Eventbus) Close() {
	self.conn.Close()
}
