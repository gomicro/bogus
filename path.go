package bogus

import ()

type Path struct {
	payload []byte
	hits    int
	status  int
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
