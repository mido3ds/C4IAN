package odmrp

import (
	"fmt"
	"time"
)

type JoinTimeout struct {
	rrep_packet      *Packet
	timeout_         Time
	delay_multiplier int
	next             *JoinTimeout
}

type MulticastTimeout struct {
	mcast_group_ int
	timeout_     Time
	next         *MulticastTimeout
}

type JoinReplyTimer struct {
	a_          *ODMRPAgent
	timers_head *JoinTimeout
	timers_tail *JoinTimeout
}

type JoinQueryTimer struct {
	a_          *ODMRPAgent
	timers_head *MulticastTimeout
	timers_tail *MulticastTimeout
}

func newJoinReplyTimer(a *ODMRPAgent) *JoinReplyTimer {
	var new_timer JoinReplyTimer
	new_timer.a_ = a
	new_timer.timers_head = nil
	new_timer.timers_tail = nil
	return &new_timer
}

func (t *JoinReplyTimer) Current() *Packet {
	if t.timers_head != nil {
		return t.timers_head.rrep_packet
	}
	return nil
}

func (t *JoinReplyTimer) CurrentTimeout() Time {
	if t.timers_head != nil {
		return t.timers_head.timeout_
	}
	return -1
}

func (t *JoinReplyTimer) EnqueueCurrent(current_time Time) bool {
	var mh, mh2 *Header
	temp := t.Dequeue()
	if t.timers_head != nil && t.timers_head.delay_multiplier == INF_NUM_RET {
		t.timers_head = t.timers_head.next
		// free temp.rrep_packet and delete temp
		temp.rrep_packet = nil
		temp = nil
		return true
	}

	if temp != nil {
		// TODO get temp.rrep_packet header
		// mh = hdr_o::access(temp.rrep_packet)
		if mh.join_reply_.count_ == 1 {
			t.Enqueue(temp.rrep_packet, current_time, temp.delay_multiplier)
			// free temp.rrep_packet and delete temp
			temp.rrep_packet = nil
			temp = nil
		} else {
			// Split reply into multiple replies because we're expecting individual responses
			for i := 0; i < mh.join_reply_.count_; i++ {
				rrep_packet := temp.rrep_packet.Copy()
				// TODO get rrep_packet header
				// mh2 = hdr_o::access(rrep_packet)

				mh2.join_reply_.pairs_[0].src_addr_ = mh.join_reply_.pairs_[i].src_addr_
				mh2.join_reply_.pairs_[0].prev_hop_addr_ = mh.join_reply_.pairs_[i].prev_hop_addr_
				mh2.join_reply_.count_ = 1

				t.Enqueue(rrep_packet, current_time, temp.delay_multiplier)
				rrep_packet = nil
			}
			// free temp.rrep_packet and delete temp
			temp.rrep_packet = nil
			temp = nil
		}
		return true
	}

	return false
}

func (t *JoinReplyTimer) Dequeue() *JoinTimeout {
	temp := t.timers_head
	var mh *Header
	// TODO get t.timers_head.rrep_packet header
	// hdr_o *mh = hdr_o::access(timers_head.rrep_packet)
	current_time := Time(time.Now().Second())

	if mh.join_reply_.count_ < 1 || t.timers_head.delay_multiplier == INF_NUM_RET {
		t.timers_head = t.timers_head.next
		temp.rrep_packet = nil
		temp = nil
		if t.timers_head != nil {
			delay := t.timers_head.timeout_ - current_time
			if delay <= 0 {
				// resched(0.001)
			} else {
				// resched(delay)
			}
		}
		return nil
	} else if t.timers_head.delay_multiplier >= MAX_NUM_RET {
		mh.join_reply_.pairs_[0].prev_hop_addr_ = -1
		mh.join_reply_.count_ = 1
		t.timers_head.delay_multiplier = INF_NUM_RET
		return temp
	}

	t.timers_head = t.timers_head.next
	return temp
}

func (t *JoinReplyTimer) Enqueue(rrep_packet *Packet, out_time Time, delay_multiplier int) {
	var new_timer *JoinTimeout
	var new_timeout Time
	new_timer = nil
	// current_time := Time(time.Now().Nanosecond() / 1000.0)
	if delay_multiplier == 0 {
		new_timeout = out_time + PACK_TIMEOUT
		delay_multiplier = 1
	} else {
		delay_multiplier += 1
		new_timeout = out_time + PACK_TIMEOUT*Time(delay_multiplier)
	}

	if t.timers_head != nil {
		t.timers_head = new(JoinTimeout)
		t.timers_head.next = nil
		t.timers_tail = t.timers_head
		t.timers_tail.next = nil
		new_timer = t.timers_tail
		// resched(new_timeout - current_time)
	} else {
		temp := t.timers_head
		/* check head first */
		if t.timers_head.timeout_ > new_timeout {
			// resched(new_timeout - current_time)
			temp = t.timers_head
			t.timers_head = new(JoinTimeout)
			t.timers_head.next = temp
			new_timer = t.timers_head
		} else {
			for temp.next != nil {
				if temp.next.timeout_ > new_timeout {
					var temp1 *JoinTimeout = nil
					temp1 = temp.next
					temp.next = new(JoinTimeout)
					temp.next.next = temp1
					new_timer = temp.next
					break
				}
				temp = temp.next
			}
			// insert at the end
			if temp.next != nil {
				temp.next = new(JoinTimeout)
				temp.next.next = nil
				new_timer = temp.next
			}

		}
	}

	new_timer.rrep_packet = rrep_packet.Copy()
	new_timer.delay_multiplier = delay_multiplier
	new_timer.timeout_ = new_timeout
}

func (t *JoinReplyTimer) LookupAndUpdate(mcast_grp int, source int, prev_hop int) bool {
	var mh *Header
	// mh2 = hdr_o::access(rrep_packet)
	temp := t.timers_head

	for temp != nil {
		// TODO get temp.rrep_packet header
		// mh = hdr_o::access(temp.rrep_packet)

		if mh.join_reply_.count_ > 0 {
			for i := 0; i < mh.join_reply_.count_; i++ {
				if mh.multicast_group_addr_ == mcast_grp && mh.join_reply_.pairs_[i].src_addr_ == source && mh.join_reply_.pairs_[i].prev_hop_addr_ == prev_hop {
					if mh.join_reply_.count_ > 1 {
						mh.join_reply_.pairs_[i].src_addr_ = mh.join_reply_.pairs_[mh.join_reply_.count_-1].src_addr_
						mh.join_reply_.pairs_[i].prev_hop_addr_ = mh.join_reply_.pairs_[mh.join_reply_.count_-1].prev_hop_addr_
					}
					mh.join_reply_.count_--
					return true
				}
			}
		}
		temp = temp.next
	}

	return false
}

func (t *JoinReplyTimer) PrintTimers() {
	temp := t.timers_head
	var mh *Header
	fmt.Printf("Timeouts: mcast grp reply_count src prev_hop\n")
	for temp != nil {
		// TODO access header of temp.rrep_packet
		// mh := hdr_o::access(temp.rrep_packet)
		fmt.Printf("%d %d %d %d\n", mh.multicast_group_addr_, mh.join_reply_.count_,
			mh.join_reply_.pairs_[0].src_addr_, mh.join_reply_.pairs_[0].prev_hop_addr_)
		temp = temp.next
	}
}

func newJoinQueryTimer(a *ODMRPAgent) *JoinQueryTimer {
	var new_timer JoinQueryTimer
	new_timer.a_ = a
	new_timer.timers_head = nil
	new_timer.timers_tail = nil
	return &new_timer
}

func (t *JoinQueryTimer) CurrentMulticastGroup() int {
	if t.timers_head != nil {
		return t.timers_head.mcast_group_
	} else {
		return -1
	}
}

func (t *JoinQueryTimer) CurrentTimeout() Time {

	if t.timers_head != nil {
		return t.timers_head.timeout_
	} else {
		return 0.0
	}
}

func (t *JoinQueryTimer) EnqueueOldGroup() {
	t.EnqueueNewGroup(t.timers_head.mcast_group_)
	t.timers_head = t.timers_head.next
}

func (t *JoinQueryTimer) EnqueueNewGroup(mcast_grp int) {
	current_time := Time(time.Now().Nanosecond() / 1000.0)

	if t.timers_head == nil {
		t.timers_head = new(MulticastTimeout)
		t.timers_head.next = nil
		t.timers_tail = t.timers_head
		t.timers_tail.next = nil
	} else {
		t.timers_tail.next = new(MulticastTimeout)
		t.timers_tail = t.timers_tail.next
		t.timers_tail.next = nil
	}

	t.timers_tail.mcast_group_ = mcast_grp
	t.timers_tail.timeout_ = current_time + JQ_REFRESH_INTERVAL
}

func (t *JoinQueryTimer) LookupGroup(mcast_grp int) bool {
	temp := t.timers_head

	for temp != nil {
		if temp.mcast_group_ == mcast_grp {
			return true
		}
		temp = temp.next
	}

	return false
}
