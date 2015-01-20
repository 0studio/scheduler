package scheduler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSchedular(t *testing.T) {
	s := InitScheduler(t, 1, 1, testSchdularFun, time.Now())
	s.Start()
	time.Sleep(3 * time.Second)

}
func testSchdularFun(s *Scheduler) {
	t := s.ID.(*testing.T)
	assert.True(t, true)
	fmt.Println("doit")
}

func TestSchedular2(t *testing.T) {
	s := InitScheduler(t, 1, 1, testSchdularFun2, time.Now())
	s.Start()
	time.Sleep(4 * time.Second)

}
func testSchdularFun2(s *Scheduler) {
	fmt.Println("it should still print this ,even after one panic")
	panic("test panic in schedular")
	t := s.ID.(*testing.T)
	assert.Fail(t, "should panic")
}
