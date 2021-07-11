package tables

import (
	"sync"
	"time"
)

const timersQueueDefaultCap = 1000

type TimersQueue struct {
	mutex    sync.Mutex
	canEnter *sync.Cond
	q        []*Timer
	nonce    uint64
}

type Timer struct {
	nonce    uint64
	t        int64
	callback func()
	gotimer  *time.Timer
	tq       *TimersQueue
}

func newTimer(d time.Duration, callback func(), nonce uint64, tq *TimersQueue) *Timer {
	return &Timer{
		nonce:    nonce,
		t:        time.Now().UnixNano() + d.Nanoseconds(),
		callback: callback,
		gotimer:  nil,
		tq:       tq,
	}
}

func (t *Timer) HasFired() bool {
	return time.Now().UnixNano() >= t.t
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already expired or been stopped.
// Note: this doesn't stop first timer
func (t *Timer) Stop() bool {
	return t.tq.stop(t)
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active, false if the timer had expired or been stopped.
// It resets the timer in either ways
func (t *Timer) Reset(d time.Duration) bool {
	return t.tq.reset(t, d)
}

func NewTimersQueue() *TimersQueue {
	return &TimersQueue{
		nonce:    0,
		q:        make([]*Timer, 0, timersQueueDefaultCap),
		canEnter: sync.NewCond(&sync.Mutex{}),
	}
}

// Add a timer that will fire after `d`, TimersQueue then will call `callback`
// in a new goroutine for each callback
func (t *TimersQueue) Add(d time.Duration, callback func()) *Timer {
	defer t.canEnter.Signal()
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.nonce += 1

	timer := newTimer(d, callback, t.nonce, t)

	for i, v := range t.q {
		if timer.t <= v.t {
			// cant insert at start
			if i == 0 {
				timer.gotimer = time.AfterFunc(d, callback)
				return timer
			}

			// insert at i
			t.q = append(t.q[:i+1], t.q[i:]...)
			t.q[i] = timer

			return timer
		}
	}

	// insert at end
	t.q = append(t.q, timer)
	return timer
}

func (tq *TimersQueue) stop(t *Timer) bool {
	if t.gotimer != nil {
		return t.gotimer.Stop()
	}

	if len(tq.q) <= 1 {
		return false
	}

	if t.HasFired() {
		return false
	}

	tq.mutex.Lock()
	defer tq.mutex.Unlock()

	for i := 0; i < len(tq.q); i++ {
		if t.nonce == tq.q[i].nonce {
			if i == 0 {
				// can't remove first entry
				return false
			}

			if i+1 == len(tq.q) {
				// remove from end
				tq.q = tq.q[:i]
			} else {
				// remove from i
				tq.q = append(tq.q[:i], tq.q[i+1:]...)
			}

			return true
		}
	}

	// not found
	return false
}

func (tq *TimersQueue) reset(t *Timer, d time.Duration) bool {
	state := tq.stop(t)
	t2 := tq.Add(d, t.callback)
	t.nonce = t2.nonce
	return state
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
