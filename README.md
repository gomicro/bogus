# Bogus
[![Build Status](https://travis-ci.org/gomicro/bogus.svg)](https://travis-ci.org/gomicro/bogus)

Bogus adds a set of helper functions around Go's httptest server.

# Usage

Setting a payload and status against a root path
```go
import "github.com/gomicro/bogus"

...

	g.Describe("Tests needing a test server", func(){
		var server *bogus.Bogus

		g.BeforeEach(func(){
			server = bogus.New()
			server.SetPayload([]byte("some return payload"))
			server.SetStatus(200)
			server.Start()
		})

		g.It("should connect to a test server", func(){
			host, port := server.HostPort()

			...

			Expect(server.Hits()).To(Equal(1))
		})
	})
```

Setting a payload and status against a specific path
```go
import "github.com/gomicro/bogus"

...

	g.Describe("Tests needing a test server", func(){
		var server *bogus.Bogus

		g.BeforeEach(func(){
			server = bogus.New()
			server.Start()
		})

		g.It("should connect to a test server", func(){
			server.AddPath("/foo/bar").
				SetPayload([]byte("some return payload")).
				SetStatus(200)
			host, port := server.HostPort()

			...

			Expect(server.Hits()).To(Equal(1))
		})
	})
```
