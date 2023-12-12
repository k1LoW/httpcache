package rfc9111

import (
	"net/http"
	"testing"
	"time"
)

func TestShared_Storable(t *testing.T) {
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
			"GET 200 -> No Store",
			&http.Request{
				Method: http.MethodGet,
			},
			&http.Response{
				StatusCode: http.StatusOK,
			},
			false,
			time.Time{},
		},
	}
	now := time.Date(2024, 12, 13, 14, 15, 16, 00, time.UTC)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s, err := NewShared()
			if err != nil {
				t.Errorf("Shared.Storable() error = %v", err)
				return
			}
			gotOK, gotExpires := s.Storable(tt.req, tt.res, now)
			if gotOK != tt.wantOK {
				t.Errorf("Shared.Storable() gotOK = %v, want %v", gotOK, tt.wantOK)
			}
			if !gotExpires.Equal(tt.wantExpires) {
				t.Errorf("Shared.Storable() gotExpires = %v, want %v", gotExpires, tt.wantExpires)
			}
		})
	}
}
