package main

import (
	"net/http"
)

type MFunc func(handler http.Handler) http.Handler

// M is an organizer function for `MFunc`s.
// All middle ware handler must call `h.next.ServeHTTP(w, r)`
// at the end of its process
func M(baseHandler http.Handler, mfuncs ...MFunc) http.Handler {
	if len(mfuncs) == 0 {
		return baseHandler
	}
	h := mfuncs[0](baseHandler)
	for i := 1; i < len(mfuncs); i++ {
		h = mfuncs[i](h)
	}
	return h
}
