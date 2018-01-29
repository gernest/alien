package alien

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseParams(t *testing.T) {
	sample := []struct {
		match, pattern, result string
	}{
		{"/hello/world", "/hello/:name", "name:world"},
		{"/let/the/bullet/fly", "/let/the/:which/:what", "which:bullet,what:fly"},
		{"/hello/to/hell.jpg", "/hello/*else", "else:to/hell.jpg"},
		{"/hello/to/hell.jpg", "/hello/to/*else", "else:hell.jpg"},
		{"/hello/to/hell.jpg", "/hello/:name/*else", "name:to,else:hell.jpg"},
		{"/everything/goes/here", "/*", "catch:everything/goes/here"},
	}

	for _, v := range sample {
		n, err := ParseParams(v.match, v.pattern)
		if err != nil {
			t.Error(err)
		}
		if n != v.result {
			t.Errorf("expected %s got %s", v.result, n)
		}
	}
}

func TestRouter(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Path))
	}
	m := New()
	_ = m.Get("/GET", h)
	_ = m.Put("/PUT", h)
	_ = m.Post("/POST", h)
	_ = m.Head("/HEAD", h)
	_ = m.Patch("/PATCH", h)
	_ = m.Options("/OPTIONS", h)
	_ = m.Connect("/CONNECT", h)
	_ = m.Trace("/TRACE", h)
	ts := httptest.NewServer(m)
	defer ts.Close()
	sample := []string{
		"GET", "POST", "PUT", "PATCH",
		"HEAD", "CONNECT", "OPTIONS", "TRACE"}
	client := &http.Client{}
	for _, v := range sample {
		req, err := http.NewRequest(v, ts.URL+"/"+v, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %d got %d %s", http.StatusOK, resp.StatusCode, req.URL.Path)
		}
		_ = resp.Body.Close()
	}
}

func TestNode(t *testing.T) {
	sample := []struct {
		path, match string
	}{
		{"/hello/:name", "/hello/world"},
		{"/hello/:namea/people/:name", "/hello/:namea/people/:name"},
		{"/home/*", "/home/alone"},
		{"/very/:name/*", "/very/complex/complicate/too/much"},
		{"/this/:is/:war", "/this/:is/:war"},
		{"/practical/", "/practical/"},
		{"/practical/joke/", "/practical/joke/"},
		{"/hell/:one/:one", "/hell/:one/:one"},
	}
	n := &node{typ: nodeRoot}
	for _, v := range sample {
		err := n.insert(v.path, &route{path: v.path})
		if err != nil {
			t.Error(err)
		}
	}
	for _, v := range sample {
		h, err := n.find(v.match)
		if err != nil {
			t.Fatal(err, v.match)
		}
		if h.path != v.path {
			t.Errorf("expected %s got %s", v.path, h.path)
		}
	}
	err := n.insert("", &route{path: ""})
	if err == nil {
		t.Error("expected an error")
	}

}

func TestRouter_mismatch(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.URL.Path))
	}
	sample := []struct {
		method, path, phony string
	}{
		{"GET", "/hello", "/"},
		{"POST", "/", "/hello"},
	}
	m := New()
	for _, v := range sample {
		_ = m.AddRoute(v.method, v.path, h)
	}

	// register unknown method
	err := m.AddRoute("CRAP", "/hell", h)
	if err == nil {
		t.Error("expected error")
	}
	ts := httptest.NewServer(m)
	defer ts.Close()
	client := &http.Client{}
	for _, v := range sample {
		req, err := http.NewRequest(v.method, ts.URL+v.phony, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected %d got %d %s", http.StatusNotFound, resp.StatusCode, req.URL.Path)
		}
		_ = resp.Body.Close()
	}
}

func TestRouter_params(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		p := GetParams(r)
		fmt.Fprint(w, p)
	}
	sample := []struct {
		path, match, params string
	}{
		{"/hello/:name", "/hello/world", "map[name:world]"},
		{"/home/*", "/home/alone", "map[catch:alone]"},
		//	{"/very/:name/*", "/very/complex/complicate/too/much", "map[name:complex catch:complicate/too/much]"},
	}
	m := New()
	for _, v := range sample {
		_ = m.Get(v.path, h)
	}

	ts := httptest.NewServer(m)
	defer ts.Close()
	client := &http.Client{}
	for _, v := range sample {
		resp, err := client.Get(ts.URL + v.match)
		if err != nil {
			t.Fatal(err)
		}
		buf := &bytes.Buffer{}
		_, _ = io.Copy(buf, resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected %d got %d ", http.StatusOK, resp.StatusCode)
		}
		if buf.String() != v.params {
			t.Errorf("expected %s got %s", v.params, buf)
		}
		_ = resp.Body.Close()
	}

}

func TestMux_Group(t *testing.T) {
	m := New()
	g := m.Group("/hello")
	_ = g.Get("/world", func(_ http.ResponseWriter, _ *http.Request) {})

	req, _ := http.NewRequest("GET", "/hello/world", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestAlienMiddlewares(t *testing.T) {
	h := func(_ http.ResponseWriter, _ *http.Request) {}

	middle := func(in http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("alien"))
			in.ServeHTTP(w, r)
		})
	}
	m := New()
	m.Use(middle)
	_ = m.Get("/", h)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, req)
	if w.Body.String() != "alien" {
		t.Errorf(" expected alien got %s ", w.Body)
	}
}
