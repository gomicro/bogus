package bogus

import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestPaths(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Paths", func() {
		var server *Bogus
		var payload = "some return payload"
		var status = http.StatusOK

		g.BeforeEach(func() {
			server = New()
			server.SetPayload([]byte(payload))
			server.SetStatus(status)
			server.Start()
		})

		g.It("should allow setting the payload for the root path", func() {
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload))
			Expect(server.Hits()).To(Equal(1))
		})

		g.It("should allow setting the return status for the root path", func() {
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(status))
			Expect(server.Hits()).To(Equal(1))
		})

		g.It("should allow adding a new path", func() {
			host, port := server.HostPort()
			server.AddPath("/foo/bar").SetStatus(http.StatusCreated)

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(status))
			Expect(server.Hits()).To(Equal(1))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(server.Hits()).To(Equal(2))
		})

		g.It("should return unique payloads per path", func() {
			host, port := server.HostPort()
			fooload := "foobar"
			server.AddPath("/foo/bar").SetPayload([]byte(fooload))

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload))
			Expect(server.Hits()).To(Equal(1))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())

			body, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(fooload))
			Expect(server.Hits()).To(Equal(2))
		})
	})
}
