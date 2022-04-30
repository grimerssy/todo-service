package cache

import (
	"sync"
)

type cfgKey string

const (
	TodoKey cfgKey = "todo"
)

type ConfigLFU struct {
	Capacities   map[cfgKey]int
	CleanupSizes map[cfgKey]int
}

type LFU struct {
	mu          sync.Mutex
	capacity    int
	cleanupSize int
	hashmap     map[interface{}]*item
	last        *item
}

func NewLFU(cfg ConfigLFU, cfgKey cfgKey) *LFU {
	capacity := cfg.Capacities[cfgKey]
	cleanupSize := cfg.CleanupSizes[cfgKey]

	switch {
	case capacity < 1:
		panic("cache capacity cannot be 0 or less")
	case cleanupSize < 1:
		panic("cache cleanup size must be 1 or more")
	case cleanupSize > capacity:
		panic("cache cleanup size cannot be bigger than capacity")
	}

	cache := &LFU{
		capacity:    capacity,
		cleanupSize: cleanupSize,
		hashmap:     make(map[interface{}]*item),
		last:        nil,
	}

	return cache
}

func (c *LFU) SetValue(key, val interface{}) {
	c.mu.Lock()

	found, ok := c.hashmap[key]
	if ok {
		found.val = val
		c.updateItem(found)
		return
	}

	item := &item{
		key: key,
		val: val,
	}
	c.addItem(item)

	c.mu.Unlock()
}

func (c *LFU) GetValue(key interface{}) interface{} {
	c.mu.Lock()

	found, ok := c.hashmap[key]
	if !ok {
		c.mu.Unlock()
		return nil
	}

	c.updateItem(found)

	c.mu.Unlock()
	return found.val
}

func (c *LFU) RemoveValue(key interface{}) {
	c.mu.Lock()

	found, ok := c.hashmap[key]
	if ok {
		c.deleteItem(found)
	}

	c.mu.Unlock()
}

type frequency struct {
	used   uint
	length uint
	first  *item
	next   *frequency
	prev   *frequency
}

type item struct {
	key  interface{}
	val  interface{}
	freq *frequency
	next *item
	prev *item
}

func (c *LFU) addItem(item *item) {
	c.hashmap[item.key] = item

	if len(c.hashmap) > c.capacity {
		for i := 0; i < c.cleanupSize; i++ {
			if c.last == nil {
				break
			}
			c.popLast()
		}
	}

	if c.last == nil {
		newF := &frequency{
			used:   1,
			length: 1,
			first:  item,
		}
		item.freq = newF

		c.last = item
		return
	}

	lastF := c.last.freq
	if lastF.used == 1 {
		item.freq = lastF

		lastF.length++
		lastF.first.append(item)
		lastF.first = item
		return
	}

	newF := &frequency{
		used:   1,
		length: 1,
		first:  item,
	}
	item.freq = newF

	newF.next = lastF
	lastF.prev = newF

	item.next = c.last
	c.last.prev = item

	c.last = item
}

func (c *LFU) updateItem(item *item) {
	oldF := item.freq
	oldF.length--

	if item.next == nil {
		oldF.first = oldF.first.prev
		newF := &frequency{
			used:   oldF.used + 1,
			length: 1,
			first:  item,
		}
		item.freq = newF

		oldF.append(newF)

		if oldF.length == 0 {
			oldF.delete()
		}
		return
	}

	if item == oldF.first {
		oldF.first = item.prev
	}

	if item == c.last && item.freq.used+1 >= item.next.freq.used {
		c.last = c.last.next
	}

	item.delete()

	if oldF.next != nil && oldF.next.used == oldF.used+1 {
		newF := oldF.next
		item.freq = newF
		newF.length++

		newF.first.append(item)
		newF.first = item

		if oldF.length == 0 {
			oldF.delete()
		}
		return
	}

	newF := &frequency{
		used:   oldF.used + 1,
		length: 1,
		first:  item,
	}
	item.freq = newF

	oldF.append(newF)

	if oldF.length == 0 {
		oldF.delete()
		return
	}

	newF.prev.first.append(item)
}

func (c *LFU) deleteItem(item *item) {
	delete(c.hashmap, item.key)
	freq := item.freq
	if item == freq.first {
		freq.first = item.prev
	}
	if item == c.last {
		c.last = item.next
	}
	freq.length--
	if freq.length == 0 {
		freq.delete()
	}
	item.delete()
}

func (c *LFU) popLast() {
	last := c.last

	delete(c.hashmap, last.key)
	lastF := last.freq
	lastF.length--
	if lastF.length == 0 {
		lastF.delete()
	}

	c.last = last.next
	last.delete()
}

func (i *item) append(newI *item) {
	newI.next = i.next
	newI.prev = i

	if i.next != nil {
		i.next.prev = newI
	}
	i.next = newI
}

func (i *item) delete() {
	if i.next != nil {
		i.next.prev = i.prev
	}
	if i.prev != nil {
		i.prev.next = i.next
	}
	i.next = nil
	i.prev = nil
}

func (f *frequency) append(newF *frequency) {
	newF.next = f.next
	newF.prev = f

	if f.next != nil {
		f.next.prev = newF
	}
	f.next = newF
}

func (f *frequency) delete() {
	if f.next != nil {
		f.next.prev = f.prev
	}
	if f.prev != nil {
		f.prev.next = f.next
	}
	f.next = nil
	f.prev = nil
}
