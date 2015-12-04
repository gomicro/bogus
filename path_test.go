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
		var host, port string

		g.BeforeEach(func() {
			server = New()
			server.SetPayload([]byte(payload))
			server.SetStatus(status)
			server.Start()
			host, port = server.HostPort()
		})

		g.It("should allow setting the payload for the root path", func() {
			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload))
			Expect(server.Hits()).To(Equal(1))
		})

		g.It("should allow setting the return status for the root path", func() {
			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(status))
			Expect(server.Hits()).To(Equal(1))
		})

		g.It("should allow adding a new path", func() {
			server.AddPath("/foo/bar").
				SetStatus(http.StatusCreated)

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
			p := "foobar"
			server.AddPath("/foo/bar").
				SetPayload([]byte(p))

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
			Expect(string(body)).To(Equal(p))
			Expect(server.Hits()).To(Equal(2))
		})

		g.It("should return the number of times a path has been hit", func() {
			p := "foobar"
			server.AddPath("/foo/bar").
				SetPayload([]byte(p))

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			Expect(server.Hits()).To(Equal(1))
			Expect(server.paths["/"].Hits()).To(Equal(1))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())

			Expect(server.Hits()).To(Equal(2))
			Expect(server.paths["/"].Hits()).To(Equal(1))
			Expect(server.paths["/foo/bar"].Hits()).To(Equal(1))
		})
		g.It("should return / for any path if only that is registered", func() {
			p := "this is / payload"
			server.SetPayload([]byte(p))

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal([]byte(p)))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())
			body, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(body).To(Equal([]byte(p)))
		})
		g.It("should respect registered paths", func() {
			payload1 := "this is / payload"
			payload2 := "this is /foo/bar payload"

			server.SetPayload([]byte(payload1))
			server.AddPath("/foo/bar").
				SetPayload([]byte(payload2))

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload1))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())

			body, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload2))
		})
		g.It("should return 404 on an unregistered path when there is more than one registration", func() {
			payload1 := "this is / payload"
			payload2 := "this is /foo/bar payload"

			server.SetPayload([]byte(payload1))
			server.AddPath("/foo/bar").
				SetPayload([]byte(payload2))

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal(payload1))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar/baz")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

			body, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Not Found"))
		})

		})
	})
}
