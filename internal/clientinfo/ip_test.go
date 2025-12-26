package clientinfo

import (
	"net/http/httptest"
	"net/netip"
	"testing"

	"git.skobk.in/skobkin/ip-detect/internal/config"
)

func TestResolveClientIP(t *testing.T) {
	t.Run("forwarded header ignored when untrusted", func(t *testing.T) {
		cfg := config.ProxyConfig{TrustForwarded: false}
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.RemoteAddr = "203.0.113.10:1234"
		req.Header.Set("X-Forwarded-For", "198.51.100.3")

		if got := resolveClientIP(req, cfg); got != "203.0.113.10" {
			t.Fatalf("expected remote IP, got %s", got)
		}
	})

	t.Run("forwarded header respected", func(t *testing.T) {
		cfg := config.ProxyConfig{TrustForwarded: true}
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.RemoteAddr = "203.0.113.10:1234"
		req.Header.Set("X-Forwarded-For", "198.51.100.3, 203.0.113.10")

		if got := resolveClientIP(req, cfg); got != "198.51.100.3" {
			t.Fatalf("expected forwarded IP, got %s", got)
		}
	})

	t.Run("trusted subnets enforced", func(t *testing.T) {
		cfg := config.ProxyConfig{
			TrustForwarded: true,
			TrustedSubnets: []netip.Prefix{mustPrefix("10.0.0.0/8")},
		}

		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.RemoteAddr = "203.0.113.10:1234"
		req.Header.Set("X-Forwarded-For", "198.51.100.3")

		if got := resolveClientIP(req, cfg); got != "203.0.113.10" {
			t.Fatalf("expected remote IP due to untrusted proxy, got %s", got)
		}
	})

	t.Run("x-real-ip fallback", func(t *testing.T) {
		cfg := config.ProxyConfig{
			TrustForwarded: true,
			TrustedSubnets: []netip.Prefix{mustPrefix("203.0.113.0/24")},
		}
		req := httptest.NewRequest("GET", "http://example.com", nil)
		req.RemoteAddr = "203.0.113.10:1234"
		req.Header.Set("X-Real-IP", "198.51.100.77")

		if got := resolveClientIP(req, cfg); got != "198.51.100.77" {
			t.Fatalf("expected X-Real-IP, got %s", got)
		}
	})
}

func mustPrefix(value string) netip.Prefix {
	prefix, err := netip.ParsePrefix(value)
	if err != nil {
		panic(err)
	}
	return prefix
}
