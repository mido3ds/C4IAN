package utils

import "sync"

type void struct{}

var voidValue void

type SetType int

const (
	NON_THREAD = iota
	THREAD
)

type Set struct {
	keys map[interface{}]void
}

type SetThread struct {
	Set
	lock sync.RWMutex
}

type SetNonThread struct {
	Set
}

type Interface interface {
	Insert(keys ...interface{})
	Remove(keys ...interface{})
	Contains(keys ...interface{}) bool
	Clear()
	Clone() Interface
	Size() int
	IsEmpty() bool
}

// Set Thread functions
func newSetNonThread() *SetNonThread {
	s := &SetNonThread{}
	s.keys = make(map[interface{}]void)

	return s
}

func newSetThread() *SetThread {
	s := &SetThread{}
	s.keys = make(map[interface{}]void)

	return s
}

func newSet(settype SetType) Interface {
	if settype == NON_THREAD {
		return newSetNonThread()
	}

	return newSetThread()
}

func (s *SetThread) Insert(keys ...interface{}) {
	// if keys empty return
	if len(keys) == 0 {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, key := range keys {
		s.keys[key] = voidValue
	}
}

func (s *SetThread) Remove(keys ...interface{}) {
	// if keys empty return
	if len(keys) == 0 {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	for _, key := range keys {
		delete(s.keys, key)
	}
}

func (s *SetThread) Contains(keys ...interface{}) bool {
	// if keys empty return true
	if len(keys) == 0 {
		return true
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	contains := true
	for _, key := range keys {
		_, contains = s.keys[key]
		if contains == false {
			return contains
		}
	}

	return contains
}

func (s *SetThread) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.keys = make(map[interface{}]void)
}

func (s *SetThread) Clone() Interface {
	cpy := newSetThread()
	for key := range s.keys {
		cpy.Insert(key)
	}

	return cpy
}

func (s *SetThread) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	len := len(s.keys)

	return len
}

func (s *SetThread) IsEmpty() bool {
	return s.Size() == 0
}

// Set Non Thread functions
func (s *SetNonThread) Insert(keys ...interface{}) {
	// if keys empty return
	if len(keys) == 0 {
		return
	}

	for _, key := range keys {
		s.keys[key] = voidValue
	}
}

func (s *SetNonThread) Remove(keys ...interface{}) {
	// if keys empty return
	if len(keys) == 0 {
		return
	}

	for _, key := range keys {
		delete(s.keys, key)
	}
}

func (s *SetNonThread) Contains(keys ...interface{}) bool {
	// if keys empty return true
	if len(keys) == 0 {
		return true
	}

	contains := true
	for _, key := range keys {
		_, contains = s.keys[key]
		if contains == false {
			return contains
		}
	}

	return contains
}

func (s *SetNonThread) Clear() {
	s.keys = make(map[interface{}]void)
}

func (s *SetNonThread) Clone() Interface {
	cpy := newSetNonThread()
	for key := range s.keys {
		cpy.Insert(key)
	}

	return cpy
}

func (s *SetNonThread) Size() int {
	len := len(s.keys)

	return len
}

func (s *SetNonThread) IsEmpty() bool {
	return s.Size() == 0
}
