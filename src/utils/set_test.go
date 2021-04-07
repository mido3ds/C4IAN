package utils

import (
	"strconv"
	"sync"
	"testing"
)

func TestSetNonThread(t *testing.T) {
	var types = []SetType{
		NON_THREAD,
		THREAD,
	}
	for _, typ := range types {
		set1 := newSet(typ)
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
		set1.Clear()
	}
}

func TestSetClear(t *testing.T) {
	var types = []SetType{
		NON_THREAD,
		THREAD,
	}
	for _, typ := range types {
		set1 := newSet(typ)
		size := 100

		for i := 0; i < size; i++ {
			set1.Insert(strconv.Itoa(i))
		}

		set1.Clear()

		// checks size = 0
		if !set1.IsEmpty() {
			t.Errorf("size = %d and it should be zero", set1.Size())
		}
		set1.Clear()
	}
}

func TestSetClone(t *testing.T) {
	var types = []SetType{
		NON_THREAD,
		THREAD,
	}
	for _, typ := range types {
		set1 := newSet(typ)
		size := 100
		for i := 0; i < size; i++ {
			set1.Insert(strconv.Itoa(i))
		}
		set2 := set1.Clone()
		for i := 0; i < size; i++ {
			key := strconv.Itoa(i)
			found1 := set1.Contains(key)
			found2 := set2.Contains(key)
			if !(found1 && found2) {
				if !found1 && !found2 {
					t.Errorf("set1 and set2 should contain: %s", key)
				} else if !found1 {
					t.Errorf("set1 should contain: %s", key)
				} else {
					t.Errorf("set2 should contain: %s", key)
				}
			}
		}
		if set1.Size() != set2.Size() {
			t.Errorf(
				"size set1 = %d and size set2 = %d and they should be equal",
				set1.Size(), set2.Size())
		}
		set1.Clear()
		set2.Clear()
	}
}

func TestSetThread(t *testing.T) {
	set1 := newSet(THREAD)

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
