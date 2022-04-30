package cache

import (
	"sync"
)

type cfgKey string

const (
	TodoKey cfgKey = "todo"
)

type ConfigLFU struct {
	Capacities map[cfgKey]int
}

type LFU struct {
	mu         sync.Mutex
	capacity   int
	hashmap    map[interface{}]*item
	last       *item
	setFunc    func(key, val interface{})
	getFunc    func(key interface{}) interface{}
	removeFunc func(key interface{})
}

func NewLFU(cfg ConfigLFU, cfgKey cfgKey) *LFU {
	cache := &LFU{
		capacity: cfg.Capacities[cfgKey],
		hashmap:  make(map[interface{}]*item),
		last:     nil,
	}
	cache.setFunc = cache.setValue
	cache.getFunc = cache.getValue
	cache.removeFunc = cache.removeValue

	if cache.capacity < 1 {
		cache.setFunc = func(_, _ interface{}) {}
		cache.getFunc = func(_ interface{}) interface{} { return nil }
		cache.removeFunc = func(_ interface{}) {}
	}

	return cache
}

func (c *LFU) SetValue(key, val interface{}) {
	c.mu.Lock()
	c.setFunc(key, val)
	c.mu.Unlock()
}

func (c *LFU) GetValue(key interface{}) interface{} {
	c.mu.Lock()
	res := c.getFunc(key)
	c.mu.Unlock()
	return res
}

func (c *LFU) RemoveValue(key interface{}) {
	c.mu.Lock()
	c.removeFunc(key)
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

func (c *LFU) setValue(key, val interface{}) {
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
}

func (c *LFU) getValue(key interface{}) interface{} {
	found, ok := c.hashmap[key]

	if !ok {
		return nil
	}

	c.updateItem(found)
	return found.val
}

func (c *LFU) removeValue(key interface{}) {
	found, ok := c.hashmap[key]
	if ok {
		c.deleteItem(found)
	}
}

func (c *LFU) addItem(item *item) {
	c.hashmap[item.key] = item

	if len(c.hashmap) > c.capacity {
		c.popLast()
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
