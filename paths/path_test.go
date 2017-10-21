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
			p := New().
				SetMethods("GET")

			Expect(len(p.methods)).To(Equal(1))
			Expect(p.methods[0]).To(Equal("GET"))

			p.SetMethods("get", "post", "put")

			Expect(len(p.methods)).To(Equal(3))
			Expect(p.methods[0]).To(Equal("GET"))
			Expect(p.methods[1]).To(Equal("POST"))
			Expect(p.methods[2]).To(Equal("PUT"))
		})

		g.It("should return true if a path has a method", func() {
			p := New().
				SetMethods("GET")

			ok := p.hasMethod("GET")
			Expect(ok).To(BeTrue())

			p.SetMethods("GET", "PUT")

			ok = p.hasMethod("PUT")
			Expect(ok).To(BeTrue())

			ok = p.hasMethod("GET")
			Expect(ok).To(BeTrue())
		})

		g.It("should return true if a path has no methods and GET is checked for", func() {
			p := New()
			ok := p.hasMethod("GET")
			Expect(ok).To(BeTrue())
		})

		g.It("should return false if a path does not have a method", func() {
			p := New().
				SetMethods("GET", "PUT")

			ok := p.hasMethod("PUT")
			Expect(ok).To(BeTrue())

			ok = p.hasMethod("GET")
			Expect(ok).To(BeTrue())

			ok = p.hasMethod("POST")
			Expect(ok).To(BeFalse())
		})

		g.It("should return false if a path has a method other than GET and GET is checked for", func() {
			p := New().
				SetMethods("PUT")

			ok := p.hasMethod("GET")
			Expect(ok).To(BeFalse())
		})
	})
}
