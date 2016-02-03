# Alien [![Coverage Status](https://coveralls.io/repos/github/gernest/alien/badge.svg?branch=master)](https://coveralls.io/github/gernest/alien?branch=master) [![Build Status](https://travis-ci.org/gernest/alien.svg?branch=master)](https://travis-ci.org/gernest/alien) [![GoDoc](https://godoc.org/github.com/gernest/alien?status.svg)](https://godoc.org/github.com/gernest/alien)

Alien is a lightweight http router( multiplexer) for Go( Golang ), made for
humans who don't like magic.

Documentation [docs](https://godoc.org/github.com/gernest/alien)

# Features

* fast ( see the benchmarks, or run them yourself)
* lightweight ( just [a single file](alien.go) read all of it in less than a minute)
* safe( designed with concurrency in mind)
* middleware support.
* routes groups
* no external dependency( only the standard library )


# Motivation
Instead of talking about how good Golang can be, I am trying to show how good Golang
can be.

Using the standard library only, following good practices and well tested code(
Over 90% coverage)

# Installation

```bash
go get github.com/gernest/alien
```

# Usage

## normal static routes

```go

package main

import (
	"log"
	"net/http"

	"github.com/gernest/alien"
)

func main() {
	m := alien.New()
	m.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	log.Fatal(http.ListenAndServe(":8090", m))
}
```

visiting your localhost at path `/` will print `hello world`

## named params

```go

package main

import (
	"log"
	"net/http"

	"github.com/gernest/alien"
)

func main() {
	m := alien.New()
	m.Get("/hello/:name", func(w http.ResponseWriter, r *http.Request) {
		p := alien.GetParams(r)
		w.Write([]byte(p.Get("name")))
	})
	log.Fatal(http.ListenAndServe(":8090", m))
}
```

visiting your localhost at path `/hello/tanzania` will print `tanzania`

## catch all params
```go
package main

import (
	"log"
	"net/http"

	"github.com/gernest/alien"
)

func main() {
	m := alien.New()
	m.Get("/hello/*name", func(w http.ResponseWriter, r *http.Request) {
		p := alien.GetParams(r)
		w.Write([]byte(p.Get("name")))
	})
	log.Fatal(http.ListenAndServe(":8090", m))
}
```

visiting your localhost at path `/hello/my/margicl/sheeplike/ship` will print 
`my/margical/sheeplike/ship`

## middlewares
Middlewares are anything that satisfy the interface
`func(http.Handler)http.Handler` . Meaning you have thousand of middlewares at
your disposal, you can use middlewares from many golang http frameworks on
alien(most support the interface).


```go

package main

import (
	"log"
	"net/http"

	"github.com/gernest/alien"
)

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello middlware"))
	})
}

func main() {
	m := alien.New()
	m.Use(middleware)
	m.Get("/", func(_ http.ResponseWriter, _ *http.Request) {
	})
	log.Fatal(http.ListenAndServe(":8090", m))
}
```

visiting your localhost at path `/` will print `hello middleware`

## groups

You can group routes

```go
package main

import (
	"log"
	"net/http"

	"github.com/gernest/alien"
)

func main() {
	m := alien.New()
	g := m.Group("/home")
	m.Use(middleware)
	g.Get("/alone", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("home alone"))
	})
	log.Fatal(http.ListenAndServe(":8090", m))
}
```

visiting your localhost at path `/home/alone` will print `home alone`

# Benchmarks
The benchmarks for alien are based on [go-hhtp-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark) for some reason I wanted to include
them in alien so anyone can benchmark for him/herself ( no more magic).

You can run all benchmarks by running the foolowing command in the root of the
alien package.

```bash
go test -bench="."
```

### memory consumption

Static  | Github  | Google+  | Parse  
-------|----------|----------|-------
132904 B |176016 B  |12720 B   |22952 B

### micro benchmarks

```bash
PASS
BenchmarkAlien_Param-4           50000      162410 ns/op    504359 B/op        9 allocs/op
BenchmarkAlien_Param5-4          20000      147724 ns/op    434577 B/op       12 allocs/op
BenchmarkAlien_Param20-4         10000      158436 ns/op    446005 B/op       28 allocs/op
BenchmarkAlien_ParamWrite-4      10000     2042458 ns/op   1162566 B/op     5033 allocs/op
BenchmarkAlien_GithubStatic-4  2000000         791 ns/op       112 B/op        3 allocs/op
BenchmarkAlien_GithubParam-4     20000      146323 ns/op    444544 B/op       10 allocs/op
BenchmarkAlien_GithubAll-4         100    14494153 ns/op  44850284 B/op     1891 allocs/op
BenchmarkAlien_GPlusStatic-4   2000000         643 ns/op        80 B/op        3 allocs/op
BenchmarkAlien_GPlusParam-4      20000      129692 ns/op    374369 B/op        9 allocs/op
BenchmarkAlien_GPlus2Params-4    10000      108297 ns/op    294492 B/op       10 allocs/op
BenchmarkAlien_GPlusAll-4        10000     6647787 ns/op  21503686 B/op      145 allocs/op
BenchmarkAlien_ParseStatic-4   2000000         697 ns/op       112 B/op        3 allocs/op
BenchmarkAlien_ParseParam-4      50000      187463 ns/op    529425 B/op       10 allocs/op
BenchmarkAlien_Parse2Params-4    20000      141359 ns/op    404469 B/op       10 allocs/op
BenchmarkAlien_ParseAll-4        10000   100911815 ns/op  40159520 B/op      230 allocs/op
BenchmarkAlien_StaticAll-4       10000      183958 ns/op     19008 B/op      471 allocs/op
ok    github.com/gernest/alien  1140.811s

```

# Contributing
Start with clicking the star button to make the author and his neighbors happy. Then fork the repository and submit a pull request for whatever change you want to be added to this project.

If you have any questions, just open an issue.

# Author
Geofrey Ernest  [@gernest](https://twitter.com/gernesti) on twitter

# Licence
MIT see [LICENCE](LICENCE)
