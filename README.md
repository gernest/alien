# Alien [![Coverage Status](https://coveralls.io/repos/github/gernest/alien/badge.svg?branch=master)](https://coveralls.io/github/gernest/alien?branch=master) [![Build Status](https://travis-ci.org/gernest/alien.svg?branch=master)](https://travis-ci.org/gernest/alien) [![GoDoc](https://godoc.org/github.com/gernest/alien?status.svg)](https://godoc.org/github.com/gernest/alien) [![Go Report Card](https://goreportcard.com/badge/github.com/gernest/alien)](https://goreportcard.com/report/github.com/gernest/alien)

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
I wanted a simple, fast, and lightweight router that has no unnecessary overhead
using the standard library only, following good practices and well tested code(
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
`func(http.Handler)http.Handler` . Meaning you have thousands of middlewares at
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
132904 B |175898 B  |12720 B   |22952 B

### micro benchmarks

```bash
PASS
BenchmarkAlien_Param-4       	 2000000	       910 ns/op	     128 B/op	       4 allocs/op
BenchmarkAlien_Param5-4      	 1000000	      1667 ns/op	     344 B/op	       8 allocs/op
BenchmarkAlien_Param20-4     	  300000	      4365 ns/op	    1664 B/op	      23 allocs/op
BenchmarkAlien_ParamWrite-4  	 1000000	      1611 ns/op	     520 B/op	       9 allocs/op
BenchmarkAlien_GithubStatic-4	 3000000	       489 ns/op	      16 B/op	       1 allocs/op
BenchmarkAlien_GithubParam-4 	 1000000	      1727 ns/op	     304 B/op	       6 allocs/op
BenchmarkAlien_GithubAll-4   	    5000	    321403 ns/op	   45601 B/op	    1043 allocs/op
BenchmarkAlien_GPlusStatic-4 	 5000000	       375 ns/op	       8 B/op	       1 allocs/op
BenchmarkAlien_GPlusParam-4  	 1000000	      1244 ns/op	     176 B/op	       5 allocs/op
BenchmarkAlien_GPlus2Params-4	 1000000	      1859 ns/op	     336 B/op	       6 allocs/op
BenchmarkAlien_GPlusAll-4    	  100000	     17470 ns/op	    2512 B/op	      62 allocs/op
BenchmarkAlien_ParseStatic-4 	 5000000	       398 ns/op	       8 B/op	       1 allocs/op
BenchmarkAlien_ParseParam-4  	 1000000	      1066 ns/op	     176 B/op	       5 allocs/op
BenchmarkAlien_Parse2Params-4	 1000000	      1364 ns/op	     256 B/op	       6 allocs/op
BenchmarkAlien_ParseAll-4    	   50000	     26921 ns/op	    3696 B/op	      93 allocs/op
BenchmarkAlien_StaticAll-4   	   10000	    116653 ns/op	    3952 B/op	     157 allocs/op
ok  	github.com/gernest/alien	27.801s

```

# Contributing
Start with clicking the star button to make the author and his neighbors happy. Then fork the repository and submit a pull request for whatever change you want to be added to this project.

If you have any questions, just open an issue.

# Author
Geofrey Ernest  [@gernesti](https://twitter.com/gernesti) on twitter

# Licence
MIT see [LICENCE](LICENCE)
