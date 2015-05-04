package scheduler

import (
	"github.com/0studio/goutils"
	"sync"
	"testing"
	"time"
)

func TestDayChange(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	now := time.Now()
	time235959 := now.Add(goutils.GetDurToNextDeadlineTimeDuration(now, 23, 59, 59))
	dc := NewDayChange()
	dc.Start(time235959)

	go func() {
		<-dc.Observe()
		wg.Done()
	}()
	go func() {
		<-dc.Observe()
		wg.Done()
	}()

	select {
	case <-time.After(time.Second * 2):
		t.FailNow()
	case <-wait(wg):
	}
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}
