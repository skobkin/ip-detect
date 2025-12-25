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
	IPAddress         string    `json:"ip_address"`
	Locale            string    `json:"locale"`
	PreferredLanguage string    `json:"preferred_language"`
	Hostname          string    `json:"hostname,omitempty"`
	UserAgent         string    `json:"user_agent,omitempty"`
	Method            string    `json:"method"`
	Path              string    `json:"path"`
	Timestamp         time.Time `json:"timestamp,omitempty"`
}

// Collect inspects the HTTP request and builds a Data snapshot.
func Collect(ctx context.Context, r *http.Request, cfg config.Config) Data {
	locale, preferred := ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	ipAddress := resolveClientIP(r, cfg.Proxy)

	data := Data{
		IPAddress:         ipAddress,
		Locale:            locale,
		PreferredLanguage: preferred,
		Method:            r.Method,
		Path:              r.URL.Path,
	}

	if cfg.Metadata.IncludeTimestamp {
		data.Timestamp = time.Now().UTC()
	}

	if cfg.Metadata.IncludeUserAgent {
		data.UserAgent = r.UserAgent()
	}

	if cfg.Resolver.EnableReverseDNS && ipAddress != "" {
		if host := reverseLookup(ctx, ipAddress, cfg.Resolver.LookupTimeout); host != "" {
			data.Hostname = host
		}
	}

	return data
}
