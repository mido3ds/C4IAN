package odmrp

import "fmt"

type CacheSlot struct {
	src int
	pid int
}

type Cache struct {
	cache      []CacheSlot
	insert_idx int
}

func (c *Cache) newCache() {
	c.cache = make([]CacheSlot, JOIN_HASH_SIZE)
	for i := 0; i < len(c.cache); i++ {
		c.cache[i].pid = -1
		c.cache[i].src = 0
	}
	c.insert_idx = 0
}

func (c *Cache) Lookup(src_addr int, pid int) IsFound {
	for i := 0; i < len(c.cache) && c.cache[i].pid != -1; i++ {
		if c.cache[i].src == src_addr && c.cache[i].pid == pid {
			return FOUND
		}
	}
	return NOT_FOUND
}

func (c *Cache) Insert(src_addr int, pid int) {
	c.cache[c.insert_idx].src = src_addr
	c.cache[c.insert_idx].pid = pid
	c.insert_idx = (c.insert_idx + 1) % len(c.cache)
}

func (c *Cache) Print() {
	if c.insert_idx < 1 {
		return
	}
	fmt.Printf("\nCache\n")
	for i := 0; i < len(c.cache) && c.cache[i].pid != -1; i++ {
		if c.cache[i].src == 1 {
			fmt.Printf("(src, pid) = (%d, %d), ", c.cache[i].src, c.cache[i].pid)
		}
	}
}
