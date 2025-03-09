package chain

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/maxwu/gotest-labels/apply"
)

// @env=local
func TestNewChain(t *testing.T) {
	t.Parallel()

	t.Run("Create an empty Chain", func(t *testing.T) {
		t.Parallel()
		chain := NewChain()
		assert.Equal(t, 0, len(chain.Links))
	})

	t.Run("Create a Chain with links", func(t *testing.T) {
		t.Parallel()
		linkA := func(h http.Handler) http.Handler {
			return h
		}
		linkB := func(h http.Handler) http.Handler {
			return h
		}
		linkC := func(h http.Handler) http.Handler {
			return h
		}
		chain := NewChain(linkA, linkB, linkC)
		assert.Equal(t, 3, len(chain.Links))
		assert.Equal(t, []LinkFunc{linkA, linkB, linkC}, chain.Links)
	})
}

// @group=unittest
// @env=ci
func TestGetHandlerWithFailingLink(t *testing.T) {
	t.Parallel()
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

	assert.Equal(t, []string{"linkA", "linkB"}, processed)
	assert.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
	body, _ := io.ReadAll(rr.Result().Body)
	assert.Contains(t, string(body), "linkB generated an error")
}

// @env=ci
func TestGetHandlerSunnyDay(t *testing.T) {
	t.Parallel()
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

	assert.Equal(t, []string{"linkA", "linkB", "linkC"}, processed)
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}
