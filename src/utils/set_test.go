package utils

import (
	"strconv"
	"sync"
	"testing"
)

func TestSetNonThread(t *testing.T) {
	set1 := NewSet(NON_THREAD)
	size := 100

	for i := 0; i < size; i++ {
		set1.Insert(strconv.Itoa(i))
	}

	// checks size = 100 & contains from 0 to 100
	if set1.Size() != size {
		size = set1.Size()
		t.Errorf("size = %d and it should be zero", size)
	}

	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		if set1.Contains(key) == false {
			t.Errorf("set should contain: %s", key)
		}
	}

	for i := 0; i < size; i++ {
		set1.Remove(strconv.Itoa(i))
	}

	// checks size = 0
	if !set1.IsEmpty() {
		t.Errorf("size = %d and it should be zero", set1.Size())
	}
}

func TestSetClear(t *testing.T) {
	set1 := NewSet(NON_THREAD)
	size := 100

	for i := 0; i < size; i++ {
		set1.Insert(strconv.Itoa(i))
	}

	set1.Clear()

	// checks size = 0
	if !set1.IsEmpty() {
		t.Errorf("size = %d and it should be zero", set1.Size())
	}
}

func TestSetThread(t *testing.T) {
	set1 := NewSet(THREAD)

	var wg sync.WaitGroup

	insert := func(key string) {
		set1.Insert(key)
		wg.Done()
	}

	remove := func(key string) {
		set1.Remove(key)
		wg.Done()
	}

	size := 100

	wg.Add(1)
	go func() {
		for i := 0; i < size; i++ {
			wg.Add(1)
			go insert(strconv.Itoa(i))
		}
		wg.Done()
	}()
	wg.Wait()

	// checks size = 100 & contains from 0 to 100
	if set1.Size() != size {
		size = set1.Size()
		t.Errorf("size = %d and it should be zero", size)
	}

	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		if set1.Contains(key) == false {
			t.Errorf("set should contain: %s", key)
		}
	}

	wg.Add(1)
	go func() {
		for i := 0; i < size; i++ {
			wg.Add(1)
			go remove(strconv.Itoa(i))
		}
		wg.Done()
	}()
	wg.Wait()

	// checks size = 0
	if !set1.IsEmpty() {
		t.Errorf("size = %d and it should be zero", set1.Size())
	}
}
