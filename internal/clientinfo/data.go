// Package clientinfo collects metadata about the client request.
package clientinfo

import (
	"context"
	"net/http"
	"time"

	"git.skobk.in/skobkin/ip-detect/internal/config"
)

// Data describes resolved request metadata that can be rendered or serialized.
type Data struct {
	IPAddress         string             `json:"ip_address"`
	Locale            *string            `json:"locale"`
	PreferredLanguage *string            `json:"preferred_language"`
	Hostname          *string            `json:"hostname"`
	UserAgent         *string            `json:"user_agent"`
	Method            string             `json:"method"`
	Path              string             `json:"path"`
	Timestamp         *time.Time         `json:"timestamp"`
	Connection        *ConnectionInfo    `json:"connection"`
	TLS               *TLSInfo           `json:"tls"`
	Proxy             *ProxyInfo         `json:"proxy"`
	Preferences       *ClientPreferences `json:"client_preferences"`
	OriginContext     *OriginContext     `json:"origin_context"`
	ClientHints       *ClientHints       `json:"ua_client_hints"`
	RequestHeaders    []HeaderEntry      `json:"request_headers"`
}

// ConnectionInfo describes the transport-level details of the request.
type ConnectionInfo struct {
	Scheme     *string `json:"scheme"`
	Protocol   *string `json:"protocol"`
	Host       *string `json:"host"`
	RemoteAddr *string `json:"remote_addr"`
}

// TLSInfo summarizes TLS session details when the request is served over HTTPS.
type TLSInfo struct {
	Version            *string `json:"version"`
	CipherSuite        *string `json:"cipher_suite"`
	ServerName         *string `json:"server_name"`
	NegotiatedProtocol *string `json:"negotiated_protocol"`
}

// ProxyInfo captures common proxy-related headers as sent by the client/proxy.
type ProxyInfo struct {
	ForwardedFor   *string `json:"forwarded_for"`
	ForwardedProto *string `json:"forwarded_proto"`
	ForwardedHost  *string `json:"forwarded_host"`
	Forwarded      *string `json:"forwarded"`
	RealIP         *string `json:"real_ip"`
	Via            *string `json:"via"`
}

// ClientPreferences describes content negotiation and privacy preferences.
type ClientPreferences struct {
	Accept                  *string `json:"accept"`
	AcceptEncoding          *string `json:"accept_encoding"`
	AcceptLanguage          *string `json:"accept_language"`
	CacheControl            *string `json:"cache_control"`
	DNT                     *string `json:"dnt"`
	UpgradeInsecureRequests *string `json:"upgrade_insecure_requests"`
	SecGPC                  *string `json:"sec_gpc"`
}

// OriginContext describes the origin/security context of the navigation.
type OriginContext struct {
	Origin       *string `json:"origin"`
	Referer      *string `json:"referer"`
	SecFetchSite *string `json:"sec_fetch_site"`
	SecFetchMode *string `json:"sec_fetch_mode"`
	SecFetchDest *string `json:"sec_fetch_dest"`
	SecFetchUser *string `json:"sec_fetch_user"`
	SecPurpose   *string `json:"sec_purpose"`
}

// ClientHints captures User-Agent Client Hints headers when present.
type ClientHints struct {
	UA              *string `json:"ua"`
	Platform        *string `json:"platform"`
	Mobile          *string `json:"mobile"`
	Model           *string `json:"model"`
	Arch            *string `json:"arch"`
	Bitness         *string `json:"bitness"`
	FullVersionList *string `json:"full_version_list"`
	PlatformVersion *string `json:"platform_version"`
}

// HeaderEntry renders a single HTTP header for HTML/JSON output.
type HeaderEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Collect inspects the HTTP request and builds a Data snapshot.
func Collect(ctx context.Context, r *http.Request, cfg config.Config) Data {
	locale, preferred := ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	ipAddress := resolveClientIP(r, cfg.Proxy)

	data := Data{
		IPAddress: ipAddress,
		Method:    r.Method,
		Path:      r.URL.Path,
	}

	data.Locale = stringPtr(locale)
	data.PreferredLanguage = stringPtr(preferred)
	if cfg.Metadata.IncludeConnection {
		data.Connection = buildConnectionInfo(r, ipAddress)
	}
	if cfg.Metadata.IncludeTLS {
		data.TLS = buildTLSInfo(r)
	}
	if cfg.Metadata.IncludeProxyDetails {
		data.Proxy = buildProxyInfo(r)
	}
	if cfg.Metadata.IncludeClientPreferences {
		data.Preferences = buildClientPreferences(r)
	}
	if cfg.Metadata.IncludeOriginContext {
		data.OriginContext = buildOriginContext(r)
	}
	if cfg.Metadata.IncludeClientHints {
		data.ClientHints = buildClientHints(r)
	}
	if cfg.Metadata.IncludeRequestHeaders {
		data.RequestHeaders = collectRequestHeaders(r)
	}

	if cfg.Metadata.IncludeTimestamp {
		now := time.Now().UTC()
		data.Timestamp = &now
	}

	if cfg.Metadata.IncludeUserAgent {
		data.UserAgent = stringPtr(r.UserAgent())
	}

	if cfg.Resolver.EnableReverseDNS && ipAddress != "" {
		if host := reverseLookup(ctx, ipAddress, cfg.Resolver.LookupTimeout); host != "" {
			data.Hostname = stringPtr(host)
		}
	}

	return data
}
