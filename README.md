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
132904 B |174096 B  |12720 B   |22952 B

### micro benchmarks

```bash
PASS
BenchmarkAlien_Param-4       	   50000	    156026 ns/op	  504352 B/op	       9 allocs/op
BenchmarkAlien_Param5-4      	   20000	    134658 ns/op	  434566 B/op	      13 allocs/op
BenchmarkAlien_Param20-4     	   10000	    140110 ns/op	  445979 B/op	      28 allocs/op
BenchmarkAlien_ParamWrite-4  	   10000	   2063635 ns/op	 1162614 B/op	    5036 allocs/op
BenchmarkAlien_GithubStatic-4	 3000000	       485 ns/op	      16 B/op	       1 allocs/op
BenchmarkAlien_GithubParam-4 	   20000	    134991 ns/op	  444530 B/op	      11 allocs/op
BenchmarkAlien_GithubAll-4   	     100	  14260420 ns/op	44846819 B/op	    1900 allocs/op
BenchmarkAlien_GPlusStatic-4 	 5000000	       376 ns/op	       8 B/op	       1 allocs/op
BenchmarkAlien_GPlusParam-4  	   30000	    174774 ns/op	  559461 B/op	      11 allocs/op
BenchmarkAlien_GPlus2Params-4	   20000	    186028 ns/op	  584668 B/op	      12 allocs/op
BenchmarkAlien_GPlusAll-4    	   10000	   5308518 ns/op	21502170 B/op	     153 allocs/op
BenchmarkAlien_ParseStatic-4 	 3000000	       396 ns/op	       8 B/op	       1 allocs/op
BenchmarkAlien_ParseParam-4  	   50000	    167199 ns/op	  529413 B/op	      10 allocs/op
BenchmarkAlien_Parse2Params-4	   20000	    128531 ns/op	  404454 B/op	      11 allocs/op
BenchmarkAlien_ParseAll-4    	   10000	   9386353 ns/op	40156535 B/op	     229 allocs/op
BenchmarkAlien_StaticAll-4   	   10000	    117060 ns/op	    3952 B/op	     157 allocs/op
ok  	github.com/gernest/alien	214.697s

```

# Contributing
Start with clicking the star button to make the author and his neighbors happy. Then fork the repository and submit a pull request for whatever change you want to be added to this project.

If you have any questions, just open an issue.

# Author
Geofrey Ernest  [@gernest](https://twitter.com/gernesti) on twitter

# Licence
MIT see [LICENCE](LICENCE)
