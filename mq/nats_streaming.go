package mq

import (
	"errors"
	"github.com/nats-io/stan.go"
)

var (
	ErrCannotConnectToNats = errors.New("Cannot connect to NATS server in Mq")
	ErrCannotConnectToMq   = errors.New("Cannot connect to Mq server by NATS conn")
)

type Mq struct {
	cluster_id string
	client_id  string
	conn       stan.Conn
	ev         *Eventbus
}

//PubAckWait
//MaxPubAcksInflight
//SetConnectionLostHandler
func NewMq(addrs []string, username, password string, cluster_id, client_id string) (*Mq, error) {
	evconn, err := NewEventbus(addrs, username, password)
	if nil != err {
		return nil, err
	}
	sc, err := stan.Connect(cluster_id, client_id,
		stan.NatsConn(evconn.GetOriginConn()),
		stan.SetConnectionLostHandler(func(c stan.Conn, err error) {
			//log.Warn("connection to mq lost.")
		}),
	)
	if nil != err {
		return nil, err
	}
	return &Mq{
		cluster_id: cluster_id,
		client_id:  client_id,
		conn:       sc,
		ev:         evconn,
	}, nil
}

func (self *Mq) Pub(subj string, msg []byte) error {
	return self.conn.Publish(subj, msg)
}

func (self *Mq) AsyncPub(subj string, msg []byte) error {
	_, err := self.conn.PublishAsync(subj, msg, func(ackNuid string, err error) {
		if err != nil {
			//TODO: handle this.
		}
	})
	return err
}

func (self *Mq) Sub(subj, durable string, cb CbSubscribe) error {
	_, err := self.conn.Subscribe(subj, func(msg *stan.Msg) {
		cb(msg.Subject, msg.Data)
	}, stan.StartWithLastReceived(), stan.DurableName(durable))
	return err
}

func (self *Mq) QSub(subj, qname, durable string, cb CbSubscribe) error {
	_, err := self.conn.QueueSubscribe(subj, qname, func(msg *stan.Msg) {
		cb(msg.Subject, msg.Data)
	}, stan.StartWithLastReceived(), stan.DurableName(durable))
	return err
}

func (self *Mq) Close() {
	self.conn.Close()
	self.ev.Close()
}
