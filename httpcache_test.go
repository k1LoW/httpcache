package httpcache

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHandlerToClientDo(t *testing.T) {
	tests := []struct {
		h    http.Handler
		req  *http.Request
		want *http.Response
	}{
		{
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "text/plain")
				_, _ = w.Write([]byte("OK"))
			}),
			&http.Request{},
			&http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"text/plain"},
				},
			},
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			got, err := HandlerToClientDo(tt.h)(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("StatusCode: got %d, want %d", got.StatusCode, tt.want.StatusCode)
			}
			if got.Header.Get("Content-Type") != tt.want.Header.Get("Content-Type") {
				t.Errorf("Content-Type: got %s, want %s", got.Header.Get("Content-Type"), tt.want.Header.Get("Content-Type"))
			}
		})
	}
}
