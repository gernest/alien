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

# Contributing
Start with clicking the star button to make the author and his neighbors happy. Then fork the repository and submit a pull request for whatever change you want to be added to this project.

If you have any questions, just open an issue.

# Author
Geofrey Ernest  [@gernesti](https://twitter.com/gernesti) on twitter

# Licence
MIT see [LICENSE](LICENSE)
