package httpcache

import (
	"net/http"
	"net/http/httptest"
	"time"
)

type Handler interface {
	Handle(req *http.Request, cachedReq *http.Request, cachedRes *http.Response, do func(*http.Request) (*http.Response, error), now time.Time) (cacheUsed bool, res *http.Response, err error)
	Storable(req *http.Request, res *http.Response, now time.Time) (ok bool, expires time.Time)
}

func HandlerToClientDo(h http.Handler) func(*http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		res := rec.Result()
		res.Header = rec.Header()
		return res, nil
	}
}
