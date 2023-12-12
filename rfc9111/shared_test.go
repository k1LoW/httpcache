package rfc9111

import (
	"net/http"
	"testing"
	"time"
)

func TestShared_Storable(t *testing.T) {
	now := time.Date(2024, 12, 13, 14, 15, 16, 00, time.UTC)

	tests := []struct {
		name        string
		req         *http.Request
		res         *http.Response
		wantOK      bool
		wantExpires time.Time
	}{
		{
			"GET 200 Cache-Control: s-maxage=10 -> +10s",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"s-maxage=10"},
				},
			},
			true,
			time.Date(2024, 12, 13, 14, 15, 26, 00, time.UTC),
		},
		{
			"GET 200 Cache-Control: max-age=15 -> +15s",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"max-age=15"},
				},
			},
			true,
			time.Date(2024, 12, 13, 14, 15, 31, 00, time.UTC),
		},
		{
			"GET 200 Expires: 2024-12-13 14:15:20",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Expires": []string{"Mon, 13 Dec 2024 14:15:20 GMT"},
				},
			},
			true,
			time.Date(2024, 12, 13, 14, 15, 20, 00, time.UTC),
		},
		{
			"GET 200 Expires: 2024-12-13 14:15:20, Date: 2024-12-13 13:15:20",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Expires": []string{"Mon, 13 Dec 2024 14:15:20 GMT"},
					"Date":    []string{"Mon, 13 Dec 2024 13:15:20 GMT"},
				},
			},
			true,
			time.Date(2024, 12, 13, 15, 15, 16, 00, time.UTC),
		},
		{
			"GET 200 Last-Modified: 2024-12-13 14:15:10, Date: 2024-12-13 14:15:20 -> +1s",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Last-Modified": []string{"Mon, 13 Dec 2024 14:15:10 GMT"},
					"Date":          []string{"Mon, 13 Dec 2024 14:15:20 GMT"},
				},
			},
			true,
			time.Date(2024, 12, 13, 14, 15, 17, 00, time.UTC),
		},
		{
			"GET 200 Last-Modified: 2024-12-13 14:15:06 -> +1s",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Last-Modified": []string{"Mon, 13 Dec 2024 14:15:06 GMT"},
				},
			},
			true,
			time.Date(2024, 12, 13, 14, 15, 17, 00, time.UTC),
		},
		{
			"GET 500 Last-Modified: 2024-12-13 14:15:06 -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusInternalServerError,
				Header: http.Header{
					"Last-Modified": []string{"Mon, 13 Dec 2024 14:15:06 GMT"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET 200 -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Date": []string{"Mon, 13 Dec 2024 14:15:10 GMT"},
				},
			},
			false,
			time.Time{},
		},
		{
			"UNUNDERSTOODMETHOD 200 -> No Store",
			&http.Request{
				Method: "UNUNDERSTOODMETHOD",
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"max-age=15"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET 100 -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusContinue,
				Header: http.Header{
					"Cache-Control": []string{"max-age=15"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET 206 -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusPartialContent,
				Header: http.Header{
					"Cache-Control": []string{"max-age=15"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET 200 Cache-Control: no-store -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"no-store"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET 200 Cache-Control: private -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"private"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET 500 Cache-Control: public, Last-Modified 2024-12-13 14:15:06 -> +1s",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Last-Modified": []string{"Mon, 13 Dec 2024 14:15:06 GMT"},
					"Cache-Control": []string{"public"},
				},
			},
			true,
			time.Date(2024, 12, 13, 14, 15, 17, 00, time.UTC),
		},
		{
			"GET 500 Cache-Control: public -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"public"},
				},
			},
			false,
			time.Time{},
		},
		{
			"GET Authorization: XXX 200 Cache-Control: max-age=15 -> No Store",
			&http.Request{
				Method: http.MethodGet,
				Header: http.Header{
					"Authorization": []string{"XXX"},
				},
			},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Cache-Control": []string{"max-age=15"},
				},
			},
			false,
			time.Time{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := NewShared()
			if err != nil {
				t.Errorf("rfc9111.Shared.Storable() error = %v", err)
				return
			}
			gotOK, gotExpires := s.Storable(tt.req, tt.res, now)
			if gotOK != tt.wantOK {
				t.Errorf("rfc9111.Shared.Storable() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if !gotExpires.Equal(tt.wantExpires) {
				t.Errorf("rfc9111.Shared.Storable() gotExpires = %v, want %v", gotExpires, tt.wantExpires)
			}
		})
	}
}
