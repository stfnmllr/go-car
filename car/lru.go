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
	"container/list"
)

type lru struct {
	keys map[interface{}]struct{}
	l    *list.List
}

func newLru() *lru {
	return &lru{
		keys: make(map[interface{}]struct{}),
		l:    list.New(),
	}
}

func (l *lru) size() int {
	return len(l.keys)
}

func (l *lru) exist(key interface{}) bool {
	_, ok := l.keys[key]
	return ok
}

func (l *lru) insertHead(key interface{}) {
	l.l.PushFront(key)
	l.keys[key] = struct{}{}
}

func (l *lru) remove(key interface{}) {
	if len(l.keys) == 0 {
		panic("lru is empty")
	}
	e := l.l.Front()
	for e.Value != key {
		e = e.Next()
	}
	l.l.Remove(e)
	delete(l.keys, key)
}

func (l *lru) removeTail() { // remove LRU entry
	if len(l.keys) == 0 {
		panic("lru is empty")
	}
	e := l.l.Back()
	key := e.Value
	l.l.Remove(e)
	delete(l.keys, key)
}
