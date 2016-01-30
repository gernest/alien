// Package alien is a lightweight http router from outer space.
package alien

import (
	"errors"
	"net/http"
	"path"
	"strings"
	"sync"
)

var (
	eof         = rune(0)
	httpMethods = struct {
		get, head, patch, put, delete, post string
		connect, options, trace             string
	}{
		"GET", "HEAD", "PATCH", "PUT", "DELETE", "POST",
		"CONNECT", "OPTIONS", "TRACE",
	}
	errRouteNotFound = errors.New("route not found")
	errBadPattern    = errors.New("bad pattern")
	errUnknownMethod = errors.New("unkown http method")
	cookieName       = "_alien"
)

type nodeType int

const (
	nodeRoot nodeType = iota
	nodeParam
	nodeNormal
	nodeCatchAll
	nodeEnd
)

type node struct {
	key      rune
	typ      nodeType
	mu       sync.RWMutex
	value    *route
	parent   *node
	children []*node
}

func (n *node) branch(key rune, val *route, typ ...nodeType) *node {
	child := &node{
		key:    key,
		value:  val,
		parent: n,
	}
	if len(typ) > 0 {
		child.typ = typ[0]
	}
	n.children = append(n.children, child)
	return child
}

func (n *node) findChild(key rune) *node {
	for _, v := range n.children {
		if v.key == key {
			return v
		}
	}
	return nil
}

func (n *node) insert(pattern string, val *route) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.typ != nodeRoot {
		return errors.New("inserting on a non root node")
	}
	size := len(pattern)
	var level *node

	for k, ch := range pattern {
		if k == 0 {
			if ch != '/' {
				return errors.New("path must start with slash ")
			}
			level = n
		}

		c := level.findChild(ch)
		switch level.typ {
		case nodeParam:
			if k < size && ch != '/' {
				continue
			}
		}
		if c != nil {
			level = c
			continue
		}
		switch ch {
		case ':':
			level = level.branch(ch, nil, nodeParam)
		case '*':
			level = level.branch(ch, nil, nodeCatchAll)
		default:
			level = level.branch(ch, nil, nodeNormal)
		}
	}
	level.branch(eof, val, nodeEnd)
	return nil
}

func (n *node) find(path string) (*route, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.typ != nodeRoot {
		return nil, errors.New("non node search")
	}
	var level *node
	var isParam bool
	for k, ch := range path {
		if k == 0 {
			level = n
		}
		c := level.findChild(ch)
		if isParam {
			if k < len(path) && ch != '/' {
				continue
			}
			isParam = false
		}
		param := level.findChild(':')
		if param != nil {
			level = param
			isParam = true
			continue
		}
		catchAll := level.findChild('*')
		if catchAll != nil {
			level = catchAll
			break
		}
		if c != nil {
			level = c
			continue
		}
		return nil, errRouteNotFound
	}
	if level != nil {
		end := level.findChild(eof)
		if end != nil {
			return end.value, nil
		}
		if slash := level.findChild('/'); slash != nil {
			end = slash.findChild(eof)
			if end != nil {
				return end.value, nil
			}
		}
	}
	return nil, errRouteNotFound
}

type route struct {
	path       string
	middleware []func(http.Handler) http.Handler
	handler    func(http.ResponseWriter, *http.Request)
}

func (r *route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var base http.Handler
	if base == nil {
		base = http.HandlerFunc(r.handler)
	}
	for _, m := range r.middleware {
		base = m(base)
	}
	base.ServeHTTP(w, req)
}

// parseParams parses params found in mateched from pattern. There are two kinds
// of params, one to capture a segment which starts with : and a nother to
// capture everything( a.k.a catch all) whis starts with *.
//
// For instance
//   pattern:="/hello/:name"
//   matched:="/hello/world"
// Will result into name:world. this function captures the named params and
// theri coreesponding values, returning them in a comma separated  string of a
// key:value nature. please see the tests for more details.
func parseParams(matched, pattern string) (result string, err error) {
	p1 := strings.Split(matched, "/")
	p2 := strings.Split(pattern, "/")
	s1 := len(p1)
	s2 := len(p2)
	if s1 < s2 {
		err = errBadPattern
		return
	}
	if strings.Contains(pattern, ":") || strings.Contains(pattern, "*") {
		for k, v := range p2 {
			if len(v) > 0 {
				switch v[0] {
				case ':':
					if len(result) == 0 {
						result = v[1:] + ":" + p1[k]
						continue
					}
					result = result + "," + v[1:] + ":" + p1[k]
				case '*':
					name := "catch"
					if k != s2-1 {
						err = errBadPattern
						return
					}
					if len(v) > 1 {
						name = v[1:]
					}
					if len(result) == 0 {
						result = name + ":" + strings.Join(p1[k:], "/")
						return
					}
					result = result + "," + name + ":" + strings.Join(p1[k:], "/")
					return
				}
			}
		}
	}
	return
}

// Params stores route params.
type Params map[string]string

// Load loads params found in src into p.
func (p Params) Load(src string) {
	s := strings.Split(src, ",")
	for _, v := range s {
		vars := strings.Split(v, ":")
		if len(vars) != 2 {
			continue
		}
		p[vars[0]] = vars[1]
	}
}

// Get returns value associated with key.
func (p Params) Get(key string) string {
	return p[key]
}

// GetParams returrns route params stored in r.
func GetParams(r *http.Request) Params {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	p := make(Params)
	p.Load(c.Value)
	return p
}

type router struct {
	get, post, patch, put, head     *node
	connect, options, trace, delete *node
}

func (r *router) addRoute(method, path string, h func(http.ResponseWriter, *http.Request), wares ...func(http.Handler) http.Handler) error {
	newRoute := &route{path: path, handler: h}
	if len(wares) > 0 {
		newRoute.middleware = append(newRoute.middleware, wares...)
	}
	switch method {
	case httpMethods.get:
		if r.get == nil {
			r.get = &node{typ: nodeRoot}
		}
		return r.get.insert(path, newRoute)
	case httpMethods.post:
		if r.post == nil {
			r.post = &node{typ: nodeRoot}
		}
		return r.post.insert(path, newRoute)
	case httpMethods.put:
		if r.put == nil {
			r.put = &node{typ: nodeRoot}
		}
		return r.put.insert(path, newRoute)
	case httpMethods.patch:
		if r.patch == nil {
			r.patch = &node{typ: nodeRoot}
		}
		return r.patch.insert(path, newRoute)
	case httpMethods.head:
		if r.head == nil {
			r.head = &node{typ: nodeRoot}
		}
		return r.head.insert(path, newRoute)
	case httpMethods.connect:
		if r.connect == nil {
			r.connect = &node{typ: nodeRoot}
		}
		return r.connect.insert(path, newRoute)
	case httpMethods.options:
		if r.options == nil {
			r.options = &node{typ: nodeRoot}
		}
		return r.options.insert(path, newRoute)
	case httpMethods.trace:
		if r.trace == nil {
			r.trace = &node{typ: nodeRoot}
		}
		return r.trace.insert(path, newRoute)
	case httpMethods.delete:
		if r.delete == nil {
			r.delete = &node{typ: nodeRoot}
		}
		return r.delete.insert(path, newRoute)
	}
	return errUnknownMethod
}

func (r *router) find(method, path string) (*route, error) {
	switch method {
	case httpMethods.get:
		if r.get != nil {
			return r.get.find(path)
		}
	case httpMethods.post:
		if r.post != nil {
			return r.post.find(path)
		}
	case httpMethods.put:
		if r.put != nil {
			return r.put.find(path)
		}
	case httpMethods.patch:
		if r.patch != nil {
			return r.patch.find(path)
		}
	case httpMethods.head:
		if r.head != nil {
			return r.head.find(path)
		}
	case httpMethods.connect:
		if r.connect != nil {
			return r.connect.find(path)
		}
	case httpMethods.options:
		if r.options != nil {
			return r.options.find(path)
		}
	case httpMethods.trace:
		if r.trace != nil {
			return r.trace.find(path)
		}
	case httpMethods.delete:
		if r.delete != nil {
			return r.delete.find(path)
		}
	}
	return nil, errRouteNotFound
}

// Mux is a http multiplexer that allow matching of http requests to the
// registered http handlers.
//
// Mux supports named parameters in urls like
//   /hello/:name
// will match
//   /hello/world
// where by inside the request passed to the handler, the param with key name and
// value world will be passed.
//
// Mux supports catch all parameters too
//   /hello/*whatever
// will match
//   /hello/world
//   /hello/world/tanzania
//   /hello/world/afica/tanzania.png
// where by inside the request passed to the handler, the param with key
// whatever will be set and valuew will be
//   world
//   world/tanzania
//   world/afica/tanzania.png
//
// If you dont specify a name in a catch all route, then the default name "catch"
// will be ussed.
type Mux struct {
	prefix     string
	middleware []func(http.Handler) http.Handler
	notFound   http.Handler
	*router
}

// New returns a new *Mux instance with default handler for mismatched routes.
func New() *Mux {
	m := &Mux{}
	m.router = &router{}
	m.notFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, errRouteNotFound.Error(), http.StatusNotFound)
	})
	return m
}

// AddRoute registers h with pattern and method. If there is a path prefix
// created via the Group method) it will be set.
func (m *Mux) AddRoute(method, pattern string, h func(http.ResponseWriter, *http.Request)) error {
	if m.prefix != "" {
		pattern = path.Join(m.prefix, pattern)
	}
	return m.addRoute(method, pattern, h, m.middleware...)
}

// Get registers h wih pattern and method GET.
func (m *Mux) Get(pattern string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.get, pattern, h)
}

// Put registers h wih pattern and method PUT.
func (m *Mux) Put(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.put, path, h)
}

// Post registers h wih pattern and method POST.
func (m *Mux) Post(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.post, path, h)
}

// Patch registers h wih pattern and method PATCH.
func (m *Mux) Patch(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.patch, path, h)
}

// Head registers h wih pattern and method HEAD.
func (m *Mux) Head(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.head, path, h)
}

// Options registers h wih pattern and method OPTIONS.
func (m *Mux) Options(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.options, path, h)
}

// Connect  registers h wih pattern and method CONNECT.
func (m *Mux) Connect(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.connect, path, h)
}

// Trace registers h wih pattern and method TRACE.
func (m *Mux) Trace(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.trace, path, h)
}

// Delete registers h wih pattern and method DELETE.
func (m *Mux) Delete(path string, h func(http.ResponseWriter, *http.Request)) error {
	return m.AddRoute(httpMethods.delete, path, h)
}

// NotFoundHandler is executed when the request route is not found.
func (m *Mux) NotFoundHandler(h http.Handler) {
	m.notFound = h
}

// ServeHTTP implements http.Handler interface. It muliplexes http requests
// against registered handlers.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	h, err := m.find(r.Method, p)
	if err != nil {
		m.notFound.ServeHTTP(w, r)
		return
	}
	params, _ := parseParams(p, h.path) // check if there is any url params
	if params != "" {
		c := &http.Cookie{Name: cookieName, Value: params}
		r.AddCookie(c)
	}
	h.ServeHTTP(w, r)
}

// Group creates a path prefix group for pattern, all routes registered using
// the returned Mux will only match if the request path starts with pattern. For
// instance .
//   m:=New()
//   home:=m.Group("/home")
//   home.Get("/alone",myHandler)
// will match
//   /home/alone
func (m *Mux) Group(pattern string) *Mux {
	return &Mux{
		prefix: pattern,
		router: m.router,
	}

}

// Use assings midlewares to the current *Mux. All routes registered by the *Mux
// after this call will have the middlewares assigned to them.
func (m *Mux) Use(middleware ...func(http.Handler) http.Handler) {
	if len(middleware) > 0 {
		m.middleware = append(m.middleware, middleware...)
	}
}
