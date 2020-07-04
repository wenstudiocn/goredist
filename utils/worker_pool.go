package utils


import (
	"time"
)

type WorkPool struct {
	goChan chan chan func(*time.Timer)
	size   int // goChan cache length
}

func NewWorkPool(size int) *WorkPool {
	return &WorkPool{
		goChan: make(chan chan func(*time.Timer), size),
		size:   size,
	}
}

// put in a task
func (wp *WorkPool) Put(f func(*time.Timer)) {
	for len(wp.goChan) > 0 {
		goRun := <-wp.goChan // read a chan from channel
		if goRun == nil {
			continue
		}
		select {
		case <-goRun: // clear？
		default:
			goRun <- f //write, execute in 'do'
			return
		}
	}
	// when len(wp.goChan) <= 0
	goRun := make(chan func(*time.Timer), 1) // create and execute
	goRun <- f
	go wp.do(goRun)
}

// in goroutine
func (wp *WorkPool) do(one chan func(*time.Timer)) {
	defer func() {
		if e := recover(); e != nil {
			//			log.Error("recover panic", zap.Any("error", e))
		}
	}()
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()

	for { // loop reading with timeout
		select {
		case <-timer.C: // timeout  = 10s
			close(one)
			f, ok := <-one
			if !ok {
				return // 关闭成功返回
			}
			f(timer)
			return
		case f, ok := <-one: // almost here
			if !ok {
				return // one 可能关闭
			}
			f(timer)
			timer.Reset(time.Second * 10)
		}
		if len(wp.goChan) < wp.size {
			wp.goChan <- one // enqueue
		}
	}
}
