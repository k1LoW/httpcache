package rfc9111

import (
	"errors"
	"strconv"
	"strings"
)

type RequestDirectives struct {
	// 5.2.1.1.  max-age
	MaxAge *uint32
	// 5.2.1.2.  max-stale
	MaxStale *uint32
	// 5.2.1.3.  min-fresh
	MinFresh *uint32
	// 5.2.1.4.  no-cache
	NoCache bool
	// 5.2.1.5.  no-store
	NoStore bool
	//5.2.1.6.  no-transform
	NoTransform bool
	// 5.2.1.7.  only-if-cached
	OnlyIfCached bool
}

type ResponseDirectives struct {
	// 5.2.2.1.  max-age
	MaxAge *uint32
	// 5.2.2.2.  must-revalidate
	MustRevalidate bool
	// 5.2.2.3.  must-understand
	MustUnderstand bool
	// 5.2.2.4.  no-cache
	NoCache bool
	// 5.2.2.5.  no-store
	NoStore bool
	// 5.2.2.6.  no-transform
	NoTransform bool
	// 5.2.2.7.  private
	Private bool
	// 5.2.2.8.  proxy-revalidate
	ProxyRevalidate bool
	// 5.2.2.9.  public
	Public bool
	// 5.2.2.10. s-maxage
	SMaxAge *uint32
}

// ParseRequestCacheControlHeader parses the Cache-Control header of a request.
func ParseRequestCacheControlHeader(headers []string) (d *RequestDirectives, errs error) {
	d = &RequestDirectives{}
	for _, h := range headers {
		tokens := strings.Split(h, ",")
		for _, t := range tokens {
			t = strings.TrimSpace(t)
			switch {
			// When there is more than one value present for a given directive (e.g., two Expires header field lines or multiple Cache-Control: max-age directives), either the first occurrence should be used or the response should be considered stale.
			case strings.HasPrefix(t, "max-age=") && d.MaxAge == nil:
				sec := strings.TrimPrefix(t, "max-age=")
				u64, err := strconv.ParseUint(sec, 10, 32)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				u32 := uint32(u64)
				d.MaxAge = &u32
			case strings.HasPrefix(t, "max-stale=") && d.MaxStale == nil:
				sec := strings.TrimPrefix(t, "max-stale=")
				u64, err := strconv.ParseUint(sec, 10, 32)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				u32 := uint32(u64)
				d.MaxAge = &u32
			case strings.HasPrefix(t, "min-fresh=") && d.MinFresh == nil:
				sec := strings.TrimPrefix(t, "min-fresh=")
				u64, err := strconv.ParseUint(sec, 10, 32)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				u32 := uint32(u64)
				d.MinFresh = &u32
			case t == "no-cache":
				d.NoCache = true
			case t == "no-store":
				d.NoStore = true
			case t == "no-transform":
				d.NoTransform = true
			case t == "only-if-cached":
				d.OnlyIfCached = true
			default:
				// A cache MUST ignore unrecognized cache directives. (https://www.rfc-editor.org/rfc/rfc9111#section-5.2.3)
			}
		}
	}
	return
}

// ParseResponseCacheControlHeader parses the Cache-Control header of a response.
func ParseResponseCacheControlHeader(headers []string) (d *ResponseDirectives, errs error) {
	d = &ResponseDirectives{}
	for _, h := range headers {
		tokens := strings.Split(h, ",")
		for _, t := range tokens {
			t = strings.TrimSpace(t)
			switch {
			// When there is more than one value present for a given directive (e.g., two Expires header field lines or multiple Cache-Control: max-age directives), either the first occurrence should be used or the response should be considered stale.
			case strings.HasPrefix(t, "max-age=") && d.MaxAge == nil:
				sec := strings.TrimPrefix(t, "max-age=")
				u64, err := strconv.ParseUint(sec, 10, 32)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				u32 := uint32(u64)
				d.MaxAge = &u32
			case t == "must-revalidate":
				d.MustRevalidate = true
			case t == "must-understand":
				d.MustUnderstand = true
			case t == "no-cache":
				d.NoCache = true
			case t == "no-store":
				d.NoStore = true
			case t == "no-transform":
				d.NoTransform = true
			case t == "private":
				d.Private = true
			case t == "proxy-revalidate":
				d.ProxyRevalidate = true
			case t == "public":
				d.Public = true
			case strings.HasPrefix(t, "s-maxage=") && d.SMaxAge == nil:
				sec := strings.TrimPrefix(t, "s-maxage=")
				u64, err := strconv.ParseUint(sec, 10, 32)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				u32 := uint32(u64)
				d.SMaxAge = &u32
			default:
				// A cache MUST ignore unrecognized cache directives. (https://www.rfc-editor.org/rfc/rfc9111#section-5.2.3)
			}
		}
	}
	return
}
