# go-car
Golang CAR cache

[![GoDoc](https://godoc.org/github.com/d024441/go-car/car?status.png)](https://godoc.org/github.com/d024441/go-car/car)

Golang CAR cache is a cache implementation using the cache replacement policy:

* Clock with adaptive replacement (CAR)

See [https://www.cse.iitd.ernet.in/~sbansal/pubs/fast04.pdf](https://www.cse.iitd.ernet.in/~sbansal/pubs/fast04.pdf) for motivation, definition
and contribution to all aspects of this implementation.

See [https://en.wikipedia.org/wiki/Cache_replacement_policies](https://en.wikipedia.org/wiki/Cache_replacement_policies) for a general discussion
of cache replacement policies.