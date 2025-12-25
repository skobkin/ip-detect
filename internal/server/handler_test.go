package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.skobk.in/skobkin/ip-detect/internal/clientinfo"
	"git.skobk.in/skobkin/ip-detect/internal/config"
)

func TestJSONEndpoint(t *testing.T) {
	handler := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	req.RemoteAddr = "198.51.100.10:1234"
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", res.Code)
	}

	if ct := res.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("expected JSON response, got %s", ct)
	}

	var payload clientinfo.Data
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	if payload.IPAddress != "198.51.100.10" {
		t.Fatalf("unexpected IP: %s", payload.IPAddress)
	}

	if payload.Locale != "en" || payload.PreferredLanguage != "en_US" {
		t.Fatalf("unexpected locale data: %+v", payload)
	}
}

func TestPlainEndpoint(t *testing.T) {
	handler := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/plain", nil)
	req.RemoteAddr = "203.0.113.42:9999"

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", res.Code)
	}

	body := strings.TrimSpace(res.Body.String())
	if body != "203.0.113.42" {
		t.Fatalf("unexpected body: %s", body)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	cfg := config.Default()
	cfg.Resolver.EnableReverseDNS = false
	cfg.Metadata.IncludeTimestamp = false

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	handler, err := newHandler(cfg, logger)
	if err != nil {
		t.Fatalf("newHandler: %v", err)
	}

	return handler
}
