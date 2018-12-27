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
	"testing"
	//	"container/list"
)

func checkInvariants(car *CAR, t *testing.T) {
	sizeT1, sizeT2, sizeB1, sizeB2 := car.t1.size(), car.t2.size(), car.b1.size(), car.b2.size()
	sizeT1T2 := sizeT1 + sizeT2
	sizeT1B1 := sizeT1 + sizeB1
	sizeT2B2 := sizeT2 + sizeB2
	sizeT1T2B1B2 := sizeT1B1 + sizeT2B2
	sizeB1B2 := sizeB1 + sizeB2

	// invariant I1
	t.Logf("I1: 0 <= %d(%d+%d) <= %d", sizeT1T2, sizeT1, sizeT2, car.c)
	if sizeT1T2 < 0 && sizeT1T2 > car.c {
		t.Fatal("invalid invariant I1 ")
	}

	// invariant I2
	t.Logf("I2: 0 <= %d(%d+%d) <= %d", sizeT1B1, sizeT1, sizeB1, car.c)
	if sizeT1B1 < 0 && sizeT1B1 > car.c {
		t.Fatal("invalid invariant I2 ")
	}

	// invariant I3
	t.Logf("I3: 0 <= %d(%d+%d) <= %d(2*%d)", sizeT2B2, sizeT1, sizeB2, 2*car.c, car.c)
	if sizeT2B2 < 0 && sizeT2B2 > 2*car.c {
		t.Fatal("invalid invariant I3 ")
	}

	// invariant I4
	t.Logf("I4: 0 <= %d(%d+%d+%d+%d) <= %d(2*%d)", sizeT1T2B1B2, sizeT1, sizeT2, sizeB1, sizeB2, 2*car.c, car.c)
	if sizeT1T2B1B2 < 0 && sizeT1T2B1B2 > 2*car.c {
		t.Fatal("invalid invariant I4 ")
	}

	// invariant I5
	if sizeT1T2 < car.c && sizeB1B2 != 0 {
		t.Fatalf("invalid invariant I5: len(b1+b2) = %d(%d+%d) - expected 0", sizeB1B2, sizeB1, sizeB2)
	}

	// invariant I6
	if sizeT1T2B1B2 >= car.c && sizeT1T2 != car.c {
		t.Fatalf("invalid invariant I6: len(t1+t2) = %d(%d+%d) - expected %d", sizeT1T2, sizeT1, sizeT2, car.c)
	}

	// invariant I7
	if car.full && sizeT1T2 != car.c {
		t.Fatalf("invalid invariant I7: cache full and len(t1+t2) = %d(%d+%d) - expected %d", sizeT1T2, sizeT1, sizeT2, car.c)
	}
}

func checkClock(clock *clock, t *testing.T) {
	// check equal size of ring and map
	if clock.head.Len() != len(clock.slots) {
		t.Fatalf("invalid clock: len(ring) = %d - len(map) = %d", clock.head.Len(), len(clock.slots))
	}

	// check equality of keys in ring and map
	e := clock.head
	for i := 0; i < clock.head.Len(); i++ {
		if _, ok := clock.slots[e.Value]; !ok {
			t.Fatalf("invalid clock: ring key %v not found in sloty", e.Value)
		}
		e = e.Next()
	}

}

func checkLru(lru *lru, t *testing.T) {
	// check equal size of list and map
	if lru.l.Len() != len(lru.keys) {
		t.Fatalf("invalid lru: len(liste) = %d - len(map) = %d", lru.l.Len(), len(lru.keys))
	}

	// check equality of keys in list and map
	e := lru.l.Front()
	for i := 0; i < lru.l.Len(); i++ {
		if _, ok := lru.keys[e.Value]; !ok {
			t.Fatalf("invalid lru: list key %v not found in keys", e.Value)
		}
		e = e.Next()
	}
}

func check(car *CAR, t *testing.T) {
	checkInvariants(car, t)
	checkClock(car.t1, t)
	checkClock(car.t2, t)
	checkLru(car.b1, t)
	checkLru(car.b2, t)
}

func TestCAR(t *testing.T) {
	car := NewCAR(16)
	check(car, t)
	car.Load(42)
	check(car, t)
}
