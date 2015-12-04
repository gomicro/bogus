package bogus

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type HitRecord struct {
	Verb  string
	Path  string
	Query url.Values
	Body  []byte
}

type Bogus struct {
	server     *httptest.Server
	hits       int
	paths      map[string]*Path
	pathsHit   chan string
	hitRecords []HitRecord
}

func New() *Bogus {
	return &Bogus{
		paths: map[string]*Path{},
	}
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
	b.pathsHit <- r.URL.Path
	b.hits++

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	b.hitRecords = append(b.hitRecords, HitRecord{r.Method, r.URL.Path, r.URL.Query(), bodyBytes})

	var path *Path
	var ok bool

	// if we have only registered the / path, use that
	// else return 404 for missing paths we've not registered
	if path, ok = b.paths[r.URL.Path]; !ok {
		if path, ok = b.paths["/"]; !ok || len(b.paths) != 1 {
			path = &Path{
				payload: []byte("Not Found"),
				status:  http.StatusNotFound,
			}
		}
	}

	path.hits++
	w.WriteHeader(path.status)
	w.Write(path.payload)
}

func (b *Bogus) Hits() int {
	return b.hits
}

func (b *Bogus) HitRecords() []HitRecord {
	return b.hitRecords
}

func (b *Bogus) HostPort() (string, string) {
	h, p, _ := net.SplitHostPort(b.server.URL[7:])
	return h, p
}

func (b *Bogus) PathHit() string {
	return <-b.pathsHit
}

func (b *Bogus) SetPayload(p []byte) {
	b.AddPath("/").SetPayload(p)
}

func (b *Bogus) SetStatus(s int) {
	b.AddPath("/").SetStatus(s)
}

func (b *Bogus) Start() {
	b.server = httptest.NewServer(http.HandlerFunc(b.HandlePaths))
	b.pathsHit = make(chan string, 1000)
}
