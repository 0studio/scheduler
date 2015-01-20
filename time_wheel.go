package scheduler

import (
	"sync"
	"time"
)

// 通过close chan  来定时通知 关注者，
// 是对 每个关注者一个timer 的优化，
// 只有一个timer ,可以有多个关注者与这个timer 进入关联

// 通过channel这种close broadcast机制，我们可以非常方便的实现一个timer，timer有一个
// channel ch，所有需要在某一个时间 “T” 收到通知的goroutine都可以尝试读该ch，当T到
// 达时候，close该ch，那么所有的goroutine都能收到该事件了。

//使用方法
// //这里我们创建了一个timingwheel，精度是1s，最大的超时等待时间为3600s
// w := timingwheel.NewTimingWheel(1 * time.Second, 3600)
// //等待10s
// <-w.After(10 * time.Second)

type TimingWheel struct {
	sync.Mutex
	interval   time.Duration // 最小精度， 多长时间check 一次
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration
	cs         []chan struct{}
	pos        int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := new(TimingWheel)
	w.interval = interval
	w.quit = make(chan struct{})
	w.pos = 0
	w.maxTimeout = time.Duration(interval * (time.Duration(buckets)))
	w.cs = make([]chan struct{}, buckets)
	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}
	w.ticker = time.NewTicker(interval)
	go w.run()

	return w
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		panic("timeout too much, over maxtimeout")
	}
	w.Lock()
	index := (w.pos + int(timeout/w.interval)) % len(w.cs)
	b := w.cs[index]
	w.Unlock()
	return b
}

func (w *TimingWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	w.Lock()
	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})
	w.pos = (w.pos + 1) % len(w.cs)
	w.Unlock()
	close(lastC)
}
