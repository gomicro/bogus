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

		g.It("should return the number of times a path has been hit", func() {
			host, port := server.HostPort()
			fooload := "foobar"
			server.AddPath("/foo/bar").SetPayload([]byte(fooload))

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
			server.SetPayload([]byte("this is / payload"))
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(bodyBytes).To(Equal([]byte("this is / payload")))

			// now try with /foo/bar not registered
			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(bodyBytes).To(Equal([]byte("this is / payload")))
		})
		g.It("should respect registered paths", func() {
			server.SetPayload([]byte("this is / payload"))
			path := server.AddPath("/foo/bar")
			path.SetPayload([]byte("this is /foo/bar payload"))
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(string(bodyBytes)).To(Equal("this is / payload"))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar")
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(string(bodyBytes)).To(Equal("this is /foo/bar payload"))

		})
		g.It("should return 404 on an unregistered path when there is more than one registration", func() {
			server.SetPayload([]byte("this is / payload"))
			path := server.AddPath("/foo/bar")
			path.SetPayload([]byte("this is /foo/bar payload"))
			host, port := server.HostPort()

			resp, err := http.Get("http://" + net.JoinHostPort(host, port))
			Expect(err).NotTo(HaveOccurred())
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(string(bodyBytes)).To(Equal("this is / payload"))

			resp, err = http.Get("http://" + net.JoinHostPort(host, port) + "/foo/bar/baz")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			resp.Body.Close()
			Expect(string(bodyBytes)).To(Equal("Not Found"))

		})
	})
}
