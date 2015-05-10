package bogus

import (
	"net"
	"net/http"
	"net/http/httptest"
)

type Bogus struct {
	server  *httptest.Server
	hits    int
	payload string
	status  int
}

func New() *Bogus {
	return &Bogus{}
}

func (b *Bogus) Close() {
	b.server.Close()
}

func (b *Bogus) Hits() int {
	return b.hits
}

func (b *Bogus) HostPort() (string, string) {
	h, p, _ := net.SplitHostPort(b.server.URL[7:])
	return h, p
}

func (b *Bogus) SetPayload(r string) {
	b.payload = r
}

func (b *Bogus) SetStatus(s int) {
	b.status = s
}

func (b *Bogus) Start() {
	b.server = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(b.status)
				b.hits++
				w.Write([]byte(b.payload))
			}))
}
