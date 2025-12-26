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

	if payload.Locale == nil || payload.PreferredLanguage == nil {
		t.Fatalf("missing locale data: %+v", payload)
	}

	if *payload.Locale != "en" || *payload.PreferredLanguage != "en_US" {
		t.Fatalf("unexpected locale data: %+v", payload)
	}

	if payload.Connection == nil || payload.Connection.RemoteAddr == nil {
		t.Fatalf("missing connection data: %+v", payload.Connection)
	}

	if *payload.Connection.RemoteAddr != "198.51.100.10" {
		t.Fatalf("unexpected connection remote addr: %s", *payload.Connection.RemoteAddr)
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

func TestJSONRespectsMetadataFlags(t *testing.T) {
	cfg := config.Default()
	cfg.Resolver.EnableReverseDNS = false
	cfg.Metadata.IncludeTimestamp = false
	cfg.Metadata.IncludeUserAgent = false
	cfg.Metadata.IncludeConnection = false
	cfg.Metadata.IncludeTLS = false
	cfg.Metadata.IncludeClientPreferences = false
	cfg.Metadata.IncludeOriginContext = false
	cfg.Metadata.IncludeClientHints = false
	cfg.Metadata.IncludeProxyDetails = false
	cfg.Metadata.IncludeRequestHeaders = false

	handler := newTestHandlerWithConfig(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	req.RemoteAddr = "203.0.113.10:8080"
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("X-Forwarded-For", "198.51.100.99")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	body := res.Body.String()
	for _, key := range []string{
		"\"connection\":null",
		"\"tls\":null",
		"\"client_preferences\":null",
		"\"origin_context\":null",
		"\"ua_client_hints\":null",
		"\"proxy\":null",
		"\"request_headers\":null",
		"\"user_agent\":null",
	} {
		if !strings.Contains(body, key) {
			t.Fatalf("expected %s in JSON, got %s", key, body)
		}
	}
}

func TestJSONIncludesHeadersWhenEnabled(t *testing.T) {
	cfg := config.Default()
	cfg.Resolver.EnableReverseDNS = false
	cfg.Metadata.IncludeRequestHeaders = true

	handler := newTestHandlerWithConfig(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	req.RemoteAddr = "203.0.113.15:8080"
	req.Header.Set("X-Test-Header", "value")

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	var payload clientinfo.Data
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	if len(payload.RequestHeaders) == 0 {
		t.Fatalf("expected request headers, got none")
	}

	found := false
	for _, header := range payload.RequestHeaders {
		if header.Key == "X-Test-Header" && header.Value == "value" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("missing X-Test-Header in response: %+v", payload.RequestHeaders)
	}
}

func TestConnectionUsesForwardedAddress(t *testing.T) {
	cfg := config.Default()
	cfg.Resolver.EnableReverseDNS = false
	cfg.Metadata.IncludeConnection = true
	cfg.Proxy.TrustForwarded = true

	handler := newTestHandlerWithConfig(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.10")

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	var payload clientinfo.Data
	if err := json.Unmarshal(res.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	if payload.Connection == nil || payload.Connection.RemoteAddr == nil {
		t.Fatalf("missing connection data: %+v", payload.Connection)
	}

	if *payload.Connection.RemoteAddr != "198.51.100.10" {
		t.Fatalf("unexpected remote addr: %s", *payload.Connection.RemoteAddr)
	}
}

func TestHTMLHidesDisabledSections(t *testing.T) {
	cfg := config.Default()
	cfg.Resolver.EnableReverseDNS = false
	cfg.Metadata.IncludeConnection = false
	cfg.Metadata.IncludeClientPreferences = false
	cfg.Metadata.IncludeOriginContext = false
	cfg.Metadata.IncludeClientHints = false
	cfg.Metadata.IncludeProxyDetails = false
	cfg.Metadata.IncludeRequestHeaders = false

	handler := newTestHandlerWithConfig(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.20:8080"
	req.Header.Set("X-Forwarded-For", "198.51.100.20")
	req.Header.Set("Via", "test-proxy")

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	body := res.Body.String()
	for _, needle := range []string{
		"<summary>Connection</summary>",
		"<summary>Client preferences</summary>",
		"<summary>Origin</summary>",
		"<summary>Client hints</summary>",
		"<summary>Request headers</summary>",
		"<h2>Proxy</h2>",
	} {
		if strings.Contains(body, needle) {
			t.Fatalf("unexpected section %s in HTML", needle)
		}
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	cfg := config.Default()
	cfg.Resolver.EnableReverseDNS = false
	cfg.Metadata.IncludeTimestamp = false

	return newTestHandlerWithConfig(t, cfg)
}

func newTestHandlerWithConfig(t *testing.T, cfg config.Config) http.Handler {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))

	handler, err := newHandler(cfg, logger)
	if err != nil {
		t.Fatalf("newHandler: %v", err)
	}

	return handler
}
