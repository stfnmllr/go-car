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

package car_test

import (
	"fmt"

	"github.com/d024441/go-car/car"
)

var text = []string{
	"Clock with adaptive replacement (CAR)",
	"Golang",
	"Hello World",
	"Cache replacement policy",
}

func loadValue(key interface{}, slotNo int) interface{} {
	return text[key.(int)]
}

func replaceValue(key interface{}, slotNo int) {
	fmt.Printf("Key %d SlotNo %d\n", key, slotNo)
}

func Example() {
	car := car.NewCAR(len(text) - 1)
	car.SetLoadValue(loadValue)
	car.SetReplaceValue(replaceValue)

	for k := range text {
		fmt.Println(car.Load(k))
	}

	// Output:
	// Clock with adaptive replacement (CAR)
	// Golang
	// Hello World
	// Key 0 SlotNo 0
	// Cache replacement policy
}
