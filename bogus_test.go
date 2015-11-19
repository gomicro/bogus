package bogus

import (
	"net"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestBogus(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Bogus Functions", func() {
		var server *Bogus

		g.BeforeEach(func() {
			server = New()
			server.Start()
		})

		g.AfterEach(func() {
			server.Close()
		})

		g.It("should return the host and port", func() {
			host, port := server.HostPort()
			Expect(host).To(Equal("127.0.0.1"))
			Expect(port).ToNot(Equal(""))
		})

		g.It("should count hits against the server", func() {
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			Expect(server.Hits()).To(Equal(1))
		})
	})
}
