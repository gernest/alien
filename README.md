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

These are results from an old 32 bit dell laptop with 2 GB of ram and 
running linux mint.

```bash
PASS
BenchmarkAlien_Param-2       	   10000	    123242 ns/op	  103802 B/op	       8 allocs/op
BenchmarkAlien_Param5-2      	   10000	    270898 ns/op	  219256 B/op	      12 allocs/op
BenchmarkAlien_Param20-2     	   10000	    646280 ns/op	  445555 B/op	      27 allocs/op
BenchmarkAlien_ParamWrite-2  	   10000	   4095764 ns/op	  711800 B/op	    5031 allocs/op
BenchmarkAlien_GithubStatic-2	 1000000	      1667 ns/op	      80 B/op	       3 allocs/op
BenchmarkAlien_GithubParam-2 	   10000	    341458 ns/op	  224242 B/op	      10 allocs/op
BenchmarkAlien_GithubAll-2   	     100	  64580632 ns/op	44826075 B/op	    1882 allocs/op
BenchmarkAlien_GPlusStatic-2 	 1000000	      1283 ns/op	      48 B/op	       3 allocs/op
BenchmarkAlien_GPlusParam-2  	   10000	    286417 ns/op	  189088 B/op	       9 allocs/op
BenchmarkAlien_GPlus2Params-2	   10000	    415275 ns/op	  294360 B/op	      10 allocs/op
BenchmarkAce_GPlusAll-2      	    5000	  18000915 ns/op	10774329 B/op	     133 allocs/op
BenchmarkAlien_ParseStatic-2 	 1000000	      1309 ns/op	      80 B/op	       3 allocs/op
BenchmarkAlien_ParseParam-2  	   10000	    164879 ns/op	  108838 B/op	       9 allocs/op
BenchmarkAlien_Parse2Params-2	   10000	    267585 ns/op	  204160 B/op	      10 allocs/op
BenchmarkAlien_ParseAll-2    	    2000	  12309365 ns/op	 8086112 B/op	     194 allocs/op
BenchmarkAlien_StaticAll-2   	    5000	    336983 ns/op	   13056 B/op	     471 allocs/op
ok  	github.com/gernest/alien	193.458s
```

# FAQ about alien

## Why name alien?
There is no magic in this package, just common sense with a juice of technology.
Aliens don't believe in magic.

## What is alien for?
Everyone, especially people who wants to understand more about Golang with a
real project( not a toy ).

## How you can  contribute to alien!

Contributions are welcome, 

* share the project with friends and family
* talk about alien on hacker news
* talk about alien on reddit
* tweet about alien
* star  alien repostory
* fork it and forget
* use alien for your new project
* buy the author a beer

# Author
Geofrey Ernest  [@gernest](https://twitter.com/gernesti) on twitter

# Licence
MIT see [LICENCE](LICENCE)
