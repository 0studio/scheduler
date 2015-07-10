package scheduler

import (
	"github.com/0studio/goutils"
	"sync"
	"time"
)

// Usage:

// 1. start the timer
// 	now := time.Now()
// dc := NewDayChange()
// dc.Start(now)

// 2 client observe the event
// <-dc.Observe()
// do something after the daychange event

type DayChange struct {
	lock         sync.Mutex
	lastTime     time.Time
	deadlineTime time.Time
	notifyChan   chan bool
}

func (dc *DayChange) Observe() <-chan bool {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	return dc.notifyChan
}
func (dc *DayChange) notify() {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	close(dc.notifyChan)
	dc.notifyChan = make(chan bool)
}

func (dc *DayChange) Start(now time.Time) {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	dc.lastTime = now
	dur := goutils.GetDurToNextDeadlineTimeDuration(dc.lastTime, 24, 0, 0)
	dc.deadlineTime = dc.lastTime.Add(dur)
	go func() {
		select {
		case <-time.After(dur + 1): // +1time.Duration make sure it is next day
			dc.notify()
			dc.Start(time.Now())
		}
	}()

}

func NewDayChange() (dc *DayChange) {
	dc = &DayChange{
		notifyChan: make(chan bool),
	}
	return
}
