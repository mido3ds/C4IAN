package tables

import (
	"testing"
	"time"
)

func TestTimersQueue(t *testing.T) {
	tim := NewTimersQueue()
	go tim.Start()

	t2 := time.Now()
	tim.Add(2*time.Second, func() {
		t.Log(1, time.Since(t2), "s")
	})
	tim.Add(1*time.Second, func() {
		t.Log(2, time.Since(t2), "s")
	})
	tim.Add(1*time.Second, func() {
		t.Log(3, time.Since(t2), "s")
	})
	tim.Add(3*time.Second, func() {
		t.Log(4, time.Since(t2), "s")
	})

	time.Sleep(5 * time.Second)
}
