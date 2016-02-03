package alien

import (
	"net/http"
	"runtime"
	"testing"
)

type testRoute struct {
	method, path string
}
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func alienHandle(_ http.ResponseWriter, _ *http.Request) {}

func alienHandleWrite(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(GetParams(r).Get("name")))
}

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := new(mockResponseWriter)
	u := r.URL
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(w, r)
	}
}

func benchRoutes(b *testing.B, router http.Handler, routes []testRoute) {
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/", nil)
	u := r.URL
	rq := u.RawQuery

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, route := range routes {
			r.Method = route.method
			r.RequestURI = route.path
			u.Path = route.path
			u.RawQuery = rq
			router.ServeHTTP(w, r)
		}
	}
}
func calcMem(name string, load func()) {

	m := new(runtime.MemStats)

	// before
	runtime.GC()
	runtime.ReadMemStats(m)
	before := m.HeapAlloc

	load()

	// after
	runtime.GC()
	runtime.ReadMemStats(m)
	after := m.HeapAlloc
	println("   "+name+":", after-before, "Bytes")
}

func loadAlien(routes []testRoute) *Mux {
	m := New()
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RequestURI))
	}
	for _, v := range routes {
		m.AddRoute(v.method, v.path, h)
	}
	return m
}

func BenchmarkAlien_Param(b *testing.B) {
	m := New()
	m.Get("/user/:name", alienHandle)
	req, _ := http.NewRequest("GET", "/user/gordon", nil)
	benchRequest(b, m, req)
}

func BenchmarkAlien_Param5(b *testing.B) {
	m := New()
	m.Get("/:a/:b/:c/:d/:e", alienHandle)
	req, _ := http.NewRequest("GET", "/test/test/test/test/test", nil)
	benchRequest(b, m, req)
}

func BenchmarkAlien_Param20(b *testing.B) {
	twentyColon := "/:a/:b/:c/:d/:e/:f/:g/:h/:i/:j/:k/:l/:m/:n/:o/:p/:q/:r/:s/:t"
	twentyRoute := "/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t"
	m := New()
	m.Get(twentyColon, alienHandle)
	req, _ := http.NewRequest("GET", twentyRoute, nil)
	benchRequest(b, m, req)
}

func BenchmarkAlien_ParamWrite(b *testing.B) {
	m := New()
	m.Get("/user/:name", alienHandleWrite)
	req, _ := http.NewRequest("GET", "/user/gordon", nil)
	benchRequest(b, m, req)
}
