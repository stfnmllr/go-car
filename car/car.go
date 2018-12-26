/*
   Golang CAR cache
   Copyright (C) 2018  Stefan Miller

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package car

import (
	"sync"
	"sync/atomic"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

const (
	refZero int32 = 0
	refOne  int32 = 1
)

type slot struct {
	_ref   int32
	slotNo int
	value  interface{}
}

func newSlot(slotNo int) *slot {
	return &slot{slotNo: slotNo}
}

func (s *slot) getRef() int32 {
	return atomic.LoadInt32(&s._ref)
}

func (s *slot) setRef(v int32) {
	atomic.StoreInt32(&s._ref, v)
}

// LoadValue is the type of a callback function to load values into the cache.
type LoadValue func(key interface{}, slotNo int) interface{}

// ReplaceValue is the type of a callback function to notify cache replacements.
type ReplaceValue func(key, value interface{}, slotNo int)

// A CAR is a cache with cache replacement policy:
// Clock with adaptive replacement (CAR).
//
type CAR struct {
	c, p   int
	t1, t2 *clock
	b1, b2 *lru
	full   bool
	slotNo int

	cbLoadValue    LoadValue
	cbReplaceValue ReplaceValue

	mu sync.RWMutex
}

// NewCAR creates a CAR cache of n slots.
func NewCAR(n int) *CAR {
	return &CAR{
		c:  n,
		t1: newClock(),
		b1: newLru(),
		t2: newClock(),
		b2: newLru(),
	}
}

// SetLoadValue sets the callback function for loading a value into the cache.
func (c *CAR) SetLoadValue(f LoadValue) {
	c.cbLoadValue = f
}

// SetReplaceValue sets the callback function for notifying cache replacements.
func (c *CAR) SetReplaceValue(f ReplaceValue) {
	c.cbReplaceValue = f
}

func (c *CAR) replaceT1() *slot {
	for {
		key, s := c.t1.removeHead()
		if s.getRef() == refZero {
			c.b1.insertHead(key)
			if c.cbReplaceValue != nil { // replace value (callback)
				c.cbReplaceValue(key, s.value, s.slotNo)
			}
			return s
		}
		s.setRef(refZero)
		c.t2.insertTail(key, s)
	}
}

func (c *CAR) replaceT2() *slot {
	for {
		key, s := c.t2.getHead()
		if s.getRef() == refZero {
			c.t2.removeHead()
			c.b2.insertHead(key)
			if c.cbReplaceValue != nil { // replace value (callback)
				c.cbReplaceValue(key, s.value, s.slotNo)
			}
			return s
		}
		s.setRef(refZero)
		c.t2.next() // make it the tail entry in T2
	}
}

func (c *CAR) hit(key interface{}) (value interface{}, ok bool) {
	if s, ok := c.t1.getSlot(key); ok { // cache hit in T1?
		s.setRef(refOne) // set reference bit to 1
		return s.value, true
	}
	if s, ok := c.t2.getSlot(key); ok { // cache hit in T2?
		s.setRef(refOne) // set reference bit to 1
		return s.value, true
	}
	return nil, false
}

func (c *CAR) miss(key interface{}) interface{} {
	inB1, inB2 := c.b1.exist(key), c.b2.exist(key)
	var s *slot

	// cache full -> replace an entry from cache
	sizeT1, sizeT2 := c.t1.size(), c.t2.size()
	if sizeT1+sizeT2 == c.c {
		c.full = true
		if sizeT1 >= max(1, c.p) {
			s = c.replaceT1()
		} else {
			s = c.replaceT2()
		}

		// cache directory replacement
		if !inB1 && !inB2 {
			sizeT1, sizeT2, sizeB1, sizeB2 := c.t1.size(), c.t2.size(), c.b1.size(), c.b2.size()
			switch {
			case sizeT1+sizeB1 == c.c:
				c.b1.removeTail() // discard LRU entry in b1
			case sizeT1+sizeT2+sizeB1+sizeB2 == 2*c.c:
				c.b2.removeTail() // discard LRU entry in b2
			}
		}
	} else {
		s = newSlot(c.slotNo)
		c.slotNo++
	}

	// cache directory miss

	if c.cbLoadValue != nil { // get value (callback)
		s.value = c.cbLoadValue(key, s.slotNo)
	} else {
		s.value = nil
	}

	switch {

	case inB1: // key is in b1
		// increase target size for t1
		div := c.b2.size() / c.b1.size()
		c.p = min(c.p+max(1, div), c.c)
		// move key from b1 to tail of t2
		c.b1.remove(key)
		c.t2.insertTail(key, s)

	case inB2: // key is in b2
		// decrease target size for t1
		div := c.b1.size() / c.b2.size()
		c.p = max(c.p-max(1, div), 0)
		// move key from b2 to tail of t2
		c.b2.remove(key)
		c.t2.insertTail(key, s)

	default: // key is not in b1 and key is not in b2
		c.t1.insertTail(key, s)

	}
	return s.value
}

// Load returns the cache value of the given cache key.
// Load is safe for concurrent use by multiple goroutines without additional locking or coordination.
func (c *CAR) Load(key interface{}) interface{} {

	// cache hit
	c.mu.RLock()
	if value, ok := c.hit(key); ok {
		c.mu.RUnlock()
		return value
	}
	c.mu.RUnlock()

	// cache miss
	c.mu.Lock()
	if value, ok := c.hit(key); ok { // test cache hit again
		c.mu.Unlock()
		return value
	}
	value := c.miss(key)
	c.mu.Unlock()
	return value
}
