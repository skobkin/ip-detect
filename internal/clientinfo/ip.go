package clientinfo

import (
	"net"
	"net/http"
	"net/netip"
	"strings"

	"git.skobk.in/skobkin/ip-detect/internal/config"
)

func resolveClientIP(r *http.Request, cfg config.ProxyConfig) string {
	remoteIP, _ := parseRemoteAddr(r.RemoteAddr)

	if cfg.TrustForwarded {
		if len(cfg.TrustedSubnets) == 0 || (remoteIP.IsValid() && ipAllowed(remoteIP, cfg.TrustedSubnets)) {
			if ip := firstForwardedIP(r.Header.Get("X-Forwarded-For")); ip.IsValid() {
				return ip.String()
			}

			if ip, ok := parseIP(r.Header.Get("X-Real-IP")); ok {
				return ip.String()
			}
		}
	}

	if remoteIP.IsValid() {
		return remoteIP.String()
	}

	return ""
}

func parseRemoteAddr(addr string) (netip.Addr, bool) {
	if addr == "" {
		return netip.Addr{}, false
	}

	if host, _, err := net.SplitHostPort(addr); err == nil {
		return parseIP(host)
	}

	return parseIP(addr)
}

func parseIP(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return netip.Addr{}, false
	}

	ip, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Addr{}, false
	}

	return ip, true
}

func firstForwardedIP(header string) netip.Addr {
	if header == "" {
		return netip.Addr{}
	}

	parts := strings.Split(header, ",")
	for _, part := range parts {
		if ip, ok := parseIP(part); ok {
			return ip
		}
	}

	return netip.Addr{}
}

func ipAllowed(ip netip.Addr, subnets []netip.Prefix) bool {
	if len(subnets) == 0 {
		return true
	}

	for _, prefix := range subnets {
		if prefix.Contains(ip) {
			return true
		}
	}

	return false
}
