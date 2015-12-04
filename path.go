package bogus

import (
	"strings"
)

type Path struct {
	payload []byte
	hits    int
	status  int
	methods []string
}

func (p *Path) Hits() int {
	return p.hits
}

func (p *Path) SetPayload(payload []byte) *Path {
	p.payload = payload
	return p
}

func (p *Path) SetStatus(s int) *Path {
	p.status = s
	return p
}

func (p *Path) SetMethods(methods ...string) *Path {
	for i, m := range methods {
		methods[i] = strings.ToUpper(m)
	}

	p.methods = methods
	return p
}
