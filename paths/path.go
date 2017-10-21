package paths

import (
	"strings"
)

// Path represents an endpoint added to a bogus server and how it should respond
type Path struct {
	Payload []byte
	Hits    int
	Status  int
	methods []string
}

// New returns a newly instantiated path object with everything initialized as
// needed.
func New() *Path {
	return &Path{}
}

// SetPayload sets the response payload for the path and returns the path for
// additional configuration
func (p *Path) SetPayload(payload []byte) *Path {
	p.Payload = payload
	return p
}

// SetStatus sets the http status for the path and returns the path for
// additional configuration
func (p *Path) SetStatus(status int) *Path {
	p.Status = status
	return p
}

// SetMethods accepts a list of methods the path should respond to
func (p *Path) SetMethods(methods ...string) *Path {
	for i, m := range methods {
		methods[i] = strings.ToUpper(m)
	}

	p.methods = methods
	return p
}

// HasMethod returns true or false based on whether a path resonds to a given
// method
func (p *Path) HasMethod(method string) bool {
	method = strings.ToUpper(method)

	// Check for configured methods first
	if len(p.methods) != 0 {
		for _, m := range p.methods {
			if m == method {
				return true
			}
		}

		return false
	}

	// If no methods configured, handle GET by default for convenience
	if method == "GET" {
		return true
	}

	return false
}
