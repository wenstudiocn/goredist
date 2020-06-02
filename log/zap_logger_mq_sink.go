package log

import (
	"code.skysarms.com/yyk/go-app-dist/mq"
)

/// zap logger sinker
/// write msg to mq

type MqSink struct {
	name string
	mq   *mq.Mq
}

func NewMqSink(name string, mq *mq.Mq) *MqSink {
	return &MqSink{
		name: name,
		mq:   mq,
	}
}

func (self *MqSink) Write(p []byte) (int, error) {
	subject := self.name
	err := self.mq.AsyncPub(subject, p)
	return len(p), err
}

func (self *MqSink) Sync() error {
	return nil
}

func (self *MqSink) Close() error {
	return nil
}
