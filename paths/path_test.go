package paths

import (
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestBogus(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Paths", func() {
		g.It("should allow setting methods for paths", func() {
			p := &Path{}

			p.SetMethods("GET")

			Expect(len(p.methods)).To(Equal(1))
			Expect(p.methods[0]).To(Equal("GET"))

			p.SetMethods("get", "post", "put")

			Expect(len(p.methods)).To(Equal(3))
			Expect(p.methods[0]).To(Equal("GET"))
			Expect(p.methods[1]).To(Equal("POST"))
			Expect(p.methods[2]).To(Equal("PUT"))
		})

		g.It("should return true if a path has a method", func() {
			p := &Path{}
			p.SetMethods("GET")

			ok := p.HasMethod("GET")
			Expect(ok).To(BeTrue())

			p.SetMethods("GET", "PUT")

			ok = p.HasMethod("PUT")
			Expect(ok).To(BeTrue())

			ok = p.HasMethod("GET")
			Expect(ok).To(BeTrue())
		})

		g.It("should return true if a path has no methods and GET is checked for", func() {
			p := &Path{}
			ok := p.HasMethod("GET")
			Expect(ok).To(BeTrue())
		})

		g.It("should return false if a path does not have a method", func() {
			p := &Path{}
			p.SetMethods("GET", "PUT")

			ok := p.HasMethod("PUT")
			Expect(ok).To(BeTrue())

			ok = p.HasMethod("GET")
			Expect(ok).To(BeTrue())

			ok = p.HasMethod("POST")
			Expect(ok).To(BeFalse())
		})

		g.It("should return false if a path has a method other than GET and GET is checked for", func() {
			p := &Path{}
			p.SetMethods("PUT")

			ok := p.HasMethod("GET")
			Expect(ok).To(BeFalse())
		})
	})
}
