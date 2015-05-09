# Bogus
Bogus adds a set of helper functions around Go's httptest server.

# Usage
```go
import "github.com/gomicro/bogus"

...

		g.Describe("Tests needing a test server", func(){
			var server *bogus.Bogus
			
			g.BeforeEach(func(){
				server = bogus.New()
				server.Payload("some return payload")
				server.Status(200)
				server.Start()
			})

			g.It("should connect to a test server", func(){
				host, port := server.HostPort()

				...

				Expect(server.Hits()).To(Equal(1))
			})
		})
```
