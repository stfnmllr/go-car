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
	"container/ring"
)

type clock struct {
	slots map[interface{}]*slot
	head  *ring.Ring
}

func newClock() *clock {
	return &clock{slots: make(map[interface{}]*slot)}
}

func (c *clock) size() int {
	return len(c.slots)
}

func (c *clock) getSlot(key interface{}) (*slot, bool) {
	if s, ok := c.slots[key]; ok {
		return s, true
	}
	return nil, false
}

func (c *clock) next() {
	if c.head != nil {
		c.head = c.head.Next()
	}
}

func (c *clock) getHead() (key interface{}, s *slot) {
	if c.head == nil {
		return nil, nil
	}
	key = c.head.Value
	return key, c.slots[key]
}

func (c *clock) removeHead() (key interface{}, s *slot) {
	switch len(c.slots) {
	case 0:
		return nil, nil
	case 1:
		key := c.head.Value
		s := c.slots[key]
		delete(c.slots, key)
		c.head = nil
		return key, s
	default:
		key := c.head.Value
		s := c.slots[key]
		delete(c.slots, key)
		r := c.head.Prev()
		r.Unlink(1) // remove head
		c.head = r.Next()
		return key, s
	}
}

func (c *clock) insertTail(key interface{}, s *slot) {
	r := ring.New(1)
	r.Value = key
	switch len(c.slots) {
	case 0:
		c.head = r
	default:
		p := c.head.Prev()
		p.Link(r)
	}
	c.slots[key] = s
}
