package tables

import (
	"testing"
	"time"
)

func TestTimersQueue(t *testing.T) {
	tim := NewTimersQueue()
	go tim.Start()

	t2 := time.Now()

	tim.Add(1*time.Second, func() {
		t.Log(2, time.Since(t2), "s")
	})

	t1 := tim.Add(2*time.Second, func() {
		t.Error("shouldnt see this log")
	})
	if !t1.Stop() {
		t.Error("couldnt stop timer")
	}

	t3 := tim.Add(2*time.Second, func() {
		t.Log(3, time.Since(t2), "s")
	})
	if !t3.Reset(5 * time.Second) {
		t.Error("couldnt reset timer")
	}

	tim.Add(3*time.Second, func() {
		t.Log(4, time.Since(t2), "s")
	})

	time.Sleep(9 * time.Second)
}
