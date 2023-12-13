package httpcache

import (
	"net/http"
	"time"
)

type Handler interface {
	Handle(req *http.Request, cachedReq *http.Request, cachedRes *http.Response, do func(*http.Request) (*http.Response, error), now time.Time) (cacheUsed bool, res *http.Response, err error)
	Storable(req *http.Request, res *http.Response, now time.Time) (ok bool, expires time.Time)
}
