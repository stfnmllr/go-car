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

/*
Package car is a cache implementation using the cache replacement policy:

- Clock with adaptive replacement (CAR)

See https://www.cse.iitd.ernet.in/~sbansal/pubs/fast04.pdf for motivation, definition
and contribution to all aspects of this package.

See https://en.wikipedia.org/wiki/Cache_replacement_policies for a general discussion
of cache replacement policies.

*/
package car
