package scheduler

import (
	"fmt"
	"runtime"
	"time"
)

const (
	FLAG_STOP_TIMER = 1
)

// 简单的定时任务
type Scheduler struct {
	input          chan int
	interval       int
	StartTime      time.Time
	schedulerTimer *time.Timer
	// tiker          *time.Ticker
	// timeChan    <-chan time.Time
	isRepeated  bool
	ID          interface{}
	fun         func(*Scheduler)
	keepRunning bool
}

// 如果interval =0，则到了expireSecs 时间之后， 不再tick,即只执行一次
func InitScheduler(ID interface{}, expireSecs int, interval int, fun func(*Scheduler), startTimeParam ...time.Time) (scheduler *Scheduler) {
	var startTime time.Time
	if len(startTimeParam) > 0 {
		startTime = startTimeParam[0]

	} else {
		startTime = time.Now()
	}

	scheduler = &Scheduler{input: make(chan int),
		ID:          ID,
		StartTime:   startTime,
		keepRunning: true,
		interval:    interval,
		fun:         fun}
	scheduler.schedulerTimer = time.NewTimer(time.Second * time.Duration(expireSecs))
	// scheduler.timeChan = scheduler.schedulerTimer.C

	// if isRepeated {
	// 	scheduler.tiker = time.NewTicker(time.Second * time.Duration(expireSecs))
	// 	scheduler.timeChan = scheduler.tiker.C
	// } else {
	// 	scheduler.schedulerTimer = time.NewTimer(time.Second * time.Duration(expireSecs))
	// 	scheduler.timeChan = scheduler.schedulerTimer.C
	// }

	return
}

func (scheduler *Scheduler) Stop() {
	if scheduler == nil {
		return
	}

	if scheduler.schedulerTimer != nil && scheduler.input != nil {
		scheduler.keepRunning = false
		select {
		case scheduler.input <- FLAG_STOP_TIMER:
			return
		case <-time.After(500 * time.Millisecond):
			return
		}
	}
}

func (scheduler *Scheduler) Start() {
	if scheduler == nil {
		return
	}

	go scheduler.run()
}

func (scheduler *Scheduler) run() {
	for scheduler.keepRunning {
		select {
		case <-scheduler.input: // 如果收到停止timer 的消息，（副本正常结束）
			break
		case <-scheduler.schedulerTimer.C: // 定时器时间到了
			scheduler.doit()
			if scheduler.interval != 0 {
				scheduler.schedulerTimer = time.NewTimer(time.Second * time.Duration(scheduler.interval))
				// break
			} else {
				break
			}
		case <-time.After(3 * time.Second): // 避免阻塞, 在当前进程里调用stop 时，会阻塞，导致Stop 直到超时里， 这里这加一个超时， 保证此for 循环 ，可以检查keepRunning 字段， 从而退出循环，以实现stop
		}
	}
	scheduler.doCancelTimer()
	// close(scheduler.input)
	scheduler.input = nil
}
func (scheduler *Scheduler) doCancelTimer() {
	if scheduler.schedulerTimer != nil {
		scheduler.schedulerTimer.Stop()
	}
	// if scheduler.tiker != nil {
	// 	scheduler.tiker.Stop()
	// }
	scheduler.schedulerTimer = nil
	// scheduler.tiker = nil
	// scheduler.timeChan = nil
}

func (scheduler *Scheduler) doit() {
	protectFunc(func() { scheduler.fun(scheduler) })
}
func (scheduler *Scheduler) IsRunning() bool {
	return scheduler.schedulerTimer != nil

}

func protectFunc(fun func()) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			for i := 0; i < 10; i++ {
				funcName, file, line, ok := runtime.Caller(i)
				if ok {
					fmt.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
				}
			}
		}
	}()
	fun()
}
