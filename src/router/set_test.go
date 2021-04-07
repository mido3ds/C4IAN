package main

import "testing"

func TestSetCreation(t *testing.T) {
	set := make(map[string]void)
	size := len(set)
	if size != 0 {
		t.Errorf("size of set = %d; want 0", size)
	}
	set["key"] = voidValue
	size = len(set)
	if size != 1 {
		t.Errorf("size of set = %d; want 1", size)
	}
	delete(set, "key")
	size = len(set)
	if size != 0 {
		t.Errorf("size of set = %d; want 0", size)
	}
}
