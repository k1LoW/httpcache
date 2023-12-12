package rfc9111

import (
	"net/http"
	"time"
)

// Shared is a shared cache that implements RFC 9111.
// The following features are not implemented
// - Private cache
// - Request directives
type Shared struct {
	understoodMethods                 []string
	understoodStatusCodes             []int
	heuristicallyCacheableStatusCodes []int
}

type SharedOption func(*Shared) error

// NewShared returns a new Shared cache handler.
func NewShared(opts ...SharedOption) (*Shared, error) {
	s := &Shared{}

	um := make([]string, len(defaultUnderstoodMethods))
	_ = copy(um, defaultUnderstoodMethods)
	s.understoodMethods = um

	us := make([]int, len(defaultUnderstoodStatusCodes))
	_ = copy(us, defaultUnderstoodStatusCodes)
	s.understoodStatusCodes = us

	hs := make([]int, len(defaultHeuristicallyCacheableStatusCodes))
	_ = copy(hs, defaultHeuristicallyCacheableStatusCodes)
	s.heuristicallyCacheableStatusCodes = hs

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (c *Shared) Storable(req *http.Request, res *http.Response, now time.Time) (bool, time.Time) {
	// 3. Storing Responses in Caches (https://www.rfc-editor.org/rfc/rfc9111#section-3)
	// - the request method is understood by the cache;
	if !contains(req.Method, c.understoodMethods) {
		return false, time.Time{}
	}

	// - the response status code is final (see https://www.rfc-editor.org/rfc/rfc9110#section-15);
	if contains(res.StatusCode, []int{
		http.StatusContinue,
		http.StatusSwitchingProtocols,
		http.StatusProcessing,
		http.StatusEarlyHints,
	}) {
		return false, time.Time{}
	}

	rescc, _ := ParseResponseCacheControlHeader(res.Header.Values("Cache-Control"))

	// - if the response status code is 206 or 304, or the must-understand cache directive (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.3) is present: the cache understands the response status code;
	if contains(res.StatusCode, []int{
		http.StatusPartialContent,
		http.StatusNotModified,
	}) || (rescc.MustUnderstand && !contains(res.StatusCode, c.understoodStatusCodes)) {
		return false, time.Time{}
	}

	// - the no-store cache directive is not present in the response (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.5);
	if rescc.NoStore {
		return false, time.Time{}
	}

	// - if the cache is shared: the private response directive is either not present or allows a shared cache to store a modified response; see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.7);
	if rescc.Private {
		return false, time.Time{}
	}

	// - if the cache is shared: the Authorization header field is not present in the request (see https://www.rfc-editor.org/rfc/rfc9111#section-11.6.2 of [HTTP]) or a response directive is present that explicitly allows shared caching (see https://www.rfc-editor.org/rfc/rfc9111#section-3.5);
	// In this specification, the following response directives have such an effect: must-revalidate (Section 5.2.2.2), public (Section 5.2.2.9), and s-maxage (Section 5.2.2.10).
	if req.Header.Get("Authorization") != "" && !rescc.MustRevalidate && !rescc.Public && rescc.SMaxAge == nil {
		return false, time.Time{}
	}

	expires := CalclateExpires(rescc, res.Header, now)
	if expires.Sub(now) <= 0 {
		return false, time.Time{}
	}

	// - the response contains at least one of the following:

	//   * a public response directive (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.9);
	if rescc.Public {
		return true, expires
	}
	//   * a private response directive, if the cache is not shared (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.7);
	// THE CACHE IS SHARED

	//   * an Expires header field (see https://www.rfc-editor.org/rfc/rfc9111#section-5.3);
	if res.Header.Get("Expires") != "" {
		return true, expires
	}
	//   * a max-age response directive (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.1);
	if rescc.MaxAge != nil {
		return true, expires
	}
	//   * if the cache is shared: an s-maxage response directive (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.2.10);
	if rescc.SMaxAge != nil {
		return true, expires
	}
	//   * a cache extension that allows it to be cached (see https://www.rfc-editor.org/rfc/rfc9111#section-5.2.3); or
	// NOT IMPLEMENTED

	//   * a status code that is defined as heuristically cacheable (see https://www.rfc-editor.org/rfc/rfc9111#section-4.2.2).
	if contains(res.StatusCode, c.heuristicallyCacheableStatusCodes) {
		return true, expires
	}

	return false, time.Time{}
}

func CalclateExpires(d *ResponseDirectives, header http.Header, now time.Time) time.Time {
	// 	4.2.1. Calculating Freshness Lifetime
	// A cache can calculate the freshness lifetime (denoted as freshness_lifetime) of a response by evaluating the following rules and using the first match:

	// - If the cache is shared and the s-maxage response directive (Section 5.2.2.10) is present, use its value, or
	if d.SMaxAge != nil {
		return now.Add(time.Duration(*d.SMaxAge) * time.Second)
	}
	// - If the max-age response directive (Section 5.2.2.1) is present, use its value, or
	if d.MaxAge != nil {
		return now.Add(time.Duration(*d.MaxAge) * time.Second)
	}
	if header.Get("Expires") != "" {
		// - If the Expires response header field (Section 5.3) is present, use its value minus the value of the Date response header field
		et, err := http.ParseTime(header.Get("Expires"))
		if err == nil {
			if header.Get("Date") != "" {
				dt, err := http.ParseTime(header.Get("Date"))
				if err == nil {
					return now.Add(et.Sub(dt))
				}
			} else {
				// (using the time the message was received if it is not present, as per Section 6.6.1 of [HTTP])
				return et // == return now.Add(et.Sub(now))
			}
		}
	}
	// Otherwise, no explicit expiration time is present in the response. A heuristic freshness lifetime might be applicable; see https://www.rfc-editor.org/rfc/rfc9111#section-4.2.2.
	if header.Get("Last-Modified") != "" {
		lt, err := http.ParseTime(header.Get("Last-Modified"))
		if err == nil {
			// If the response has a Last-Modified header field (Section 8.8.2 of [HTTP]), caches are encouraged to use a heuristic expiration value that is no more than some fraction of the interval since that time. A typical setting of this fraction might be 10%.
			if header.Get("Date") != "" {
				dt, err := http.ParseTime(header.Get("Date"))
				if err == nil {
					return now.Add(dt.Sub(lt) / 10)
				}
			} else {
				return now.Add(now.Sub(lt) / 10)
			}
		}
	}

	return now
}

func contains[T comparable](v T, vv []T) bool {
	for _, vvv := range vv {
		if vvv == v {
			return true
		}
	}
	return false
}
