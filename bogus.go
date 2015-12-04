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

	var status int
	var payload []byte

	// if we have only registered the / path, use that
	if path, ok := b.paths["/"]; ok && len(b.paths) == 1 {
		path.hits++
		status = path.status
		payload = path.payload
	} else {
		// if we've registered the given path, let's use it.
		// else if we've not registered a path, return 404
		if path, ok := b.paths[r.URL.Path]; ok {
			path.hits++
			status = path.status
			payload = path.payload
		} else {
			status = http.StatusNotFound
			payload = []byte("Not Found")
		}
	}

	w.WriteHeader(status)
	w.Write(payload)
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
