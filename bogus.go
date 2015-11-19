package bogus

import (
	"net"
	"net/http"
	"net/http/httptest"
)

type Bogus struct {
	server *httptest.Server
	hits   int
	paths  map[string]*Path
}

func New() *Bogus {
	return &Bogus{paths: map[string]*Path{
		"/": &Path{},
	}}
}

func (b *Bogus) AddPath(path string) *Path {
	if _, ok := b.paths[path]; !ok {
		b.paths[path] = &Path{}
	}

	return b.paths[path]
}

func (b *Bogus) Close() {
	b.server.Close()
}

func (b *Bogus) HandlePaths(w http.ResponseWriter, r *http.Request) {
	if path, ok := b.paths[r.URL.Path]; ok {
		w.WriteHeader(path.status)
		b.hits++
		path.hits++
		w.Write(path.payload)
	}
}

func (b *Bogus) Hits() int {
	return b.hits
}

func (b *Bogus) HostPort() (string, string) {
	h, p, _ := net.SplitHostPort(b.server.URL[7:])
	return h, p
}

func (b *Bogus) SetPayload(p []byte) {
	path := b.paths["/"]
	if path != nil {
		path.payload = p
	}
}

func (b *Bogus) SetStatus(s int) {
	path := b.paths["/"]
	if path != nil {
		path.status = s
	}
}

func (b *Bogus) Start() {
	b.server = httptest.NewServer(http.HandlerFunc(b.HandlePaths))
}
