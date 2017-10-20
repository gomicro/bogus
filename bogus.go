package bogus

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// HitRecord represents a recording of information from a single hit againstr
// the bogus server
type HitRecord struct {
	Verb  string
	Path  string
	Query url.Values
	Body  []byte
}

// Bogus represents a test server
type Bogus struct {
	server     *httptest.Server
	hits       int
	paths      map[string]*Path
	hitRecords []HitRecord
}

// New returns a newly intitated bogus server
func New() *Bogus {
	b := &Bogus{
		paths: map[string]*Path{},
	}
	b.server = httptest.NewServer(http.HandlerFunc(b.HandlePaths))

	return b
}

// AddPath adds a new path to the bogus server handler and returns the new path
// for further configuration
func (b *Bogus) AddPath(path string) *Path {
	if _, ok := b.paths[path]; !ok {
		b.paths[path] = &Path{}
	}

	return b.paths[path]
}

// Close calls the close method for the underlying httptest server
func (b *Bogus) Close() {
	b.server.Close()
}

// HandlePaths implements the http handler interface and decides how to respond
// based on the paths configured
func (b *Bogus) HandlePaths(w http.ResponseWriter, r *http.Request) {
	b.hits++

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	b.hitRecords = append(b.hitRecords, HitRecord{r.Method, r.URL.Path, r.URL.Query(), bodyBytes})

	path, ok := b.paths[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	if path.hasMethod(r.Method) {
		var payload []byte
		var status int

		switch r.Method {
		case "POST":
			payload = bodyBytes
			status = http.StatusAccepted
		case "PUT":
			payload = bodyBytes
			status = http.StatusCreated
		case "DELETE":
			payload = []byte("")
			status = http.StatusNoContent
		case "", "GET":
			payload = path.payload
			status = http.StatusOK
		}

		// Prefer set payload and status over default
		if path.payload != nil {
			payload = path.payload
		}
		if path.status != 0 {
			status = path.status
		}

		path.hits++
		w.WriteHeader(status)
		w.Write(payload)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(""))
	return
}

// Hits returns the total number of hits seen against the bogus server
func (b *Bogus) Hits() int {
	return b.hits
}

// HitRecords returns a slice of the hit records recorded for inspection
func (b *Bogus) HitRecords() []HitRecord {
	return b.hitRecords
}

// HostPort returns the host and port number of the bogus server
func (b *Bogus) HostPort() (string, string) {
	h, p, _ := net.SplitHostPort(b.server.URL[7:])
	return h, p
}
