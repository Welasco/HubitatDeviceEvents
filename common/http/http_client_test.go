package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPClientCancelsHeadersAndBodyReads(t *testing.T) {
	for _, stage := range []string{"headers", "body"} {
		t.Run(stage, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if stage == "body" {
					w.Header().Set("Content-Length", "100")
					w.WriteHeader(http.StatusOK)
					io.WriteString(w, "partial")
					w.(http.Flusher).Flush()
				}
				<-r.Context().Done()
			}))
			defer server.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			start := time.Now()
			_, err := Http_client(ctx, server.URL)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected deadline error, got %v", err)
			}
			if elapsed := time.Since(start); elapsed > time.Second {
				t.Errorf("request did not honor deadline: %s", elapsed)
			}
		})
	}
}

func TestHTTPClientReportsTruncatedBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		io.WriteString(w, "partial")
	}))
	defer server.Close()
	body, err := Http_client(context.Background(), server.URL)
	if !errors.Is(err, io.ErrUnexpectedEOF) || body != "" {
		t.Fatalf("expected truncated-body error without partial success, got %q, %v", body, err)
	}
}
