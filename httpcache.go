package httpcache

import (
	"net/http"
	"time"
)

type Handler interface { //nostyle:ifacenames
	Storable(req *http.Request, res *http.Response, now time.Time) (ok bool, expires time.Time)
}
