package utils

import (
	"github.com/wenstudiocn/goredist/e"
	"sync"
	"time"
)

// A map with an expire, which let the expired key removed automatically.
// it has not a good performance in this version implementation.

type emItem struct {
	val       interface{}
	expiresAt int64
}

type ExpiredMap struct {
	sync.Mutex
	d     map[string]*emItem
	index map[int64][]string // TODO: for future usage to optimize
	chQ   <-chan struct{}
}

func NewExpiredMap(ch <-chan struct{}) *ExpiredMap {
	em := &ExpiredMap{
		d:   make(map[string]*emItem),
		chQ: ch,
	}
	return em
}

func (self *ExpiredMap) worker_in_background() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-self.chQ:
			return
		case <-ticker.C:
			self.remove_expired_items()
		}
	}
}

func (self *ExpiredMap) remove_expired_items() {
	now := time.Now().Unix()
	self.Lock()
	defer self.Unlock()

	for k, v := range self.d {
		if now <= v.expiresAt {
			delete(self.d, k)
		}
	}
}

//@elapse: positive value(+) normally. negative value presents a permenent key.
func (self *ExpiredMap) Put(key string, val interface{}, elapse time.Duration) error {
	// we dont check elapse if negative
	expiresAt := time.Now().Add(elapse).Unix()
	item := &emItem{
		val:       val,
		expiresAt: expiresAt,
	}
	self.Lock()
	self.d[key] = item
	self.Unlock()
	return nil
}

func (self *ExpiredMap) Expire(key string, elapse time.Duration) error {
	expiresAt := time.Now().Add(elapse).Unix()
	self.Lock()
	defer self.Unlock()

	item, ok := self.d[key]
	if !ok {
		return e.ErrKeyNotExists
	}
	item.expiresAt = expiresAt
	return nil
}

func (self *ExpiredMap) Get(key string) (interface{}, error) {
	self.Lock()
	defer self.Unlock()

	item, ok := self.d[key]
	if !ok {
		return nil, e.ErrKeyNotExists
	}
	return item.val, nil
}

func (self *ExpiredMap) Remove(key string) error {
	self.Lock()
	defer self.Unlock()
	_, ok := self.d[key]
	if !ok {
		return e.ErrKeyNotExists
	}
	delete(self.d, key)
	return nil
}

// @cb: break foreach if return true
func (self *ExpiredMap) Foreach(cb func(string, interface{}) bool) {
	self.Lock()
	defer self.Unlock()

	for k, v := range self.d {
		if cb(k, v.val) {
			break
		}
	}
}

func (self *ExpiredMap) Clear() {
	self.Lock()
	self.d = make(map[string]*emItem)
	self.Unlock()
}
