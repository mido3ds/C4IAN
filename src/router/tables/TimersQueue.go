package tables

import (
	"sync"
	"time"
)

const timersQueueDefaultCap = 1000

type TimersQueue struct {
	mutex    sync.Mutex
	canEnter *sync.Cond
	q        []timer
}

type timer struct {
	t        int64
	callback func()
}

func newTimer(d time.Duration, callback func()) timer {
	return timer{
		t:        time.Now().UnixNano() + d.Nanoseconds(),
		callback: callback,
	}
}

func NewTimersQueue() *TimersQueue {
	return &TimersQueue{
		q:        make([]timer, 0, timersQueueDefaultCap),
		canEnter: sync.NewCond(&sync.Mutex{}),
	}
}

// Add a timer that will fire after `d`, TimersQueue then will call `callback`
// in a new goroutine for each callback
func (t *TimersQueue) Add(d time.Duration, callback func()) {
	defer t.canEnter.Signal()
	t.mutex.Lock()
	defer t.mutex.Unlock()

	timer := newTimer(d, callback)

	for i, v := range t.q {
		if timer.t <= v.t {
			// cant insert at start
			if i == 0 {
				time.AfterFunc(d, callback)
				return
			}

			// insert at i
			t.q = append(t.q[:i+1], t.q[i:]...)
			t.q[i] = timer

			return
		}
	}

	// insert at end
	t.q = append(t.q, timer)
}

func nanosToDuration(nanos int64) time.Duration {
	return time.Duration((nanos - time.Now().UnixNano()) * int64(time.Nanosecond))
}

func (t *TimersQueue) Start() {
	for {
		// wait for entry
		t.canEnter.L.Lock()
		for len(t.q) == 0 {
			t.canEnter.Wait()
		}
		t.canEnter.L.Unlock()

		// sleep
		time.Sleep(nanosToDuration(t.q[0].t))
		go t.q[0].callback()

		// remove
		t.mutex.Lock()
		t.q = t.q[1:]
		t.mutex.Unlock()
	}
}
