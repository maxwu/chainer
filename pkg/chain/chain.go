package chain

import "slices"

import "net/http"

type LinkFunc func(http.Handler) http.Handler

type Chain struct {
	Links []LinkFunc
}

// NewChain creates an http handler chain from LinkFunc typed links
func NewChain(links ...LinkFunc) Chain {
	return Chain{
		slices.Clone(links),
	}
}

// GetHandler returns the final assembled http.Handler from its links
func (c Chain) GetHandler() http.Handler {
	h := http.Handler(
		http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
	)
	for i := len(c.Links) - 1; i >= 0; i-- {
		h = c.Links[i](h)
	}
	return h
}
