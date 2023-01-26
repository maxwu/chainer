package chain

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var tr *testing.T

var _ = Describe("Test Chain", func() {
	It("Verify the Chain stops processing on failed Link", func() {
		processed := []string{}
		linkA := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				processed = append(processed, "linkA")
				h.ServeHTTP(w, r)
			})
		}

		linkB := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				processed = append(processed, "linkB")
				http.Error(w, "linkB generated an error", http.StatusInternalServerError)
			})
		}

		linkC := func(_ http.Handler) http.Handler {
			return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				processed = append(processed, "linkC")
			})
		}

		h := NewChain(linkA, linkB, linkC).GetHandler()
		req, _ := http.NewRequest("GET", "/path", nil)
		rr := httptest.NewRecorder()

		// handle an incoming request
		h.ServeHTTP(rr, req)

		Expect(processed).To(Equal([]string{"linkA", "linkB"}))
		Expect(rr.Result().StatusCode).To(Equal(http.StatusInternalServerError))
		body, _ := io.ReadAll(rr.Result().Body)
		Expect(string(body)).To(ContainSubstring("linkB generated an error"))
	})

	It("Sunnyday scenario", func() {
		processed := []string{}
		linkA := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				processed = append(processed, "linkA")
				h.ServeHTTP(w, r)
			})
		}

		linkB := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				processed = append(processed, "linkB")
				h.ServeHTTP(w, r)
			})
		}

		linkC := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				processed = append(processed, "linkC")
				w.WriteHeader(http.StatusOK)
				h.ServeHTTP(w, r)
			})
		}

		h := NewChain(linkA, linkB, linkC).GetHandler()
		req, _ := http.NewRequest("POST", "/doesnt-matter", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		Expect(processed).To(Equal([]string{"linkA", "linkB", "linkC"}))
		Expect(rr.Result().StatusCode).To(Equal(http.StatusOK))
	})
})

func TestRunner(t *testing.T) {
	tr = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chain Test Suite")
}
