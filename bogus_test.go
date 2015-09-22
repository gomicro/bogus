package bogus

import (
	"io/ioutil"
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
		var payload = "some return payload"
		var status = 200

		g.BeforeEach(func() {
			server = New()
			server.SetPayload([]byte(payload))
			server.SetStatus(status)
			server.Start()
		})

		g.It("should return the host and port", func() {
			host, port := server.HostPort()
			Expect(host).To(Equal("127.0.0.1"))
			Expect(port).ToNot(Equal(""))
		})

		g.It("should allow setting the payload", func() {
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload))
			Expect(server.Hits()).To(Equal(1))
		})

		g.It("should allow setting the return status", func() {
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(status))
			Expect(server.Hits()).To(Equal(1))
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
