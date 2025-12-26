package clientinfo

import (
	"crypto/tls"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func buildConnectionInfo(r *http.Request, resolvedIP string) *ConnectionInfo {
	remoteAddr := stringPtr(resolvedIP)
	if remoteAddr == nil {
		if addr, ok := parseRemoteAddr(r.RemoteAddr); ok {
			remoteAddr = stringPtr(addr.String())
		}
	}

	info := ConnectionInfo{
		Scheme:     stringPtr(detectScheme(r)),
		Protocol:   stringPtr(r.Proto),
		Host:       stringPtr(r.Host),
		RemoteAddr: remoteAddr,
	}

	if !hasPtr(info.Scheme, info.Protocol, info.Host, info.RemoteAddr) {
		return nil
	}

	return &info
}

func buildTLSInfo(r *http.Request) *TLSInfo {
	if r.TLS == nil {
		return nil
	}

	info := TLSInfo{
		Version:            stringPtr(tls.VersionName(r.TLS.Version)),
		CipherSuite:        stringPtr(tls.CipherSuiteName(r.TLS.CipherSuite)),
		ServerName:         stringPtr(r.TLS.ServerName),
		NegotiatedProtocol: stringPtr(r.TLS.NegotiatedProtocol),
	}

	if !hasPtr(info.Version, info.CipherSuite, info.ServerName, info.NegotiatedProtocol) {
		return nil
	}

	return &info
}

func buildProxyInfo(r *http.Request) *ProxyInfo {
	info := ProxyInfo{
		ForwardedFor:   stringPtr(headerValue(r, "X-Forwarded-For")),
		ForwardedProto: stringPtr(headerValue(r, "X-Forwarded-Proto")),
		ForwardedHost:  stringPtr(headerValue(r, "X-Forwarded-Host")),
		Forwarded:      stringPtr(headerValue(r, "Forwarded")),
		RealIP:         stringPtr(headerValue(r, "X-Real-IP")),
		Via:            stringPtr(headerValue(r, "Via")),
	}

	if !hasPtr(info.ForwardedFor, info.ForwardedProto, info.ForwardedHost, info.Forwarded, info.RealIP, info.Via) {
		return nil
	}

	return &info
}

func buildClientPreferences(r *http.Request) *ClientPreferences {
	info := ClientPreferences{
		Accept:                  stringPtr(headerValue(r, "Accept")),
		AcceptEncoding:          stringPtr(headerValue(r, "Accept-Encoding")),
		AcceptLanguage:          stringPtr(headerValue(r, "Accept-Language")),
		CacheControl:            stringPtr(headerValue(r, "Cache-Control")),
		DNT:                     stringPtr(headerValue(r, "DNT")),
		UpgradeInsecureRequests: stringPtr(headerValue(r, "Upgrade-Insecure-Requests")),
		SecGPC:                  stringPtr(headerValue(r, "Sec-GPC")),
	}

	if !hasPtr(info.Accept, info.AcceptEncoding, info.AcceptLanguage, info.CacheControl, info.DNT, info.UpgradeInsecureRequests, info.SecGPC) {
		return nil
	}

	return &info
}

func buildOriginContext(r *http.Request) *OriginContext {
	info := OriginContext{
		Origin:       stringPtr(headerValue(r, "Origin")),
		Referer:      stringPtr(headerValue(r, "Referer")),
		SecFetchSite: stringPtr(headerValue(r, "Sec-Fetch-Site")),
		SecFetchMode: stringPtr(headerValue(r, "Sec-Fetch-Mode")),
		SecFetchDest: stringPtr(headerValue(r, "Sec-Fetch-Dest")),
		SecFetchUser: stringPtr(headerValue(r, "Sec-Fetch-User")),
		SecPurpose:   stringPtr(headerValue(r, "Sec-Purpose")),
	}

	if !hasPtr(info.Origin, info.Referer, info.SecFetchSite, info.SecFetchMode, info.SecFetchDest, info.SecFetchUser, info.SecPurpose) {
		return nil
	}

	return &info
}

func buildClientHints(r *http.Request) *ClientHints {
	info := ClientHints{
		UA:              stringPtr(headerValue(r, "Sec-CH-UA")),
		Platform:        stringPtr(headerValue(r, "Sec-CH-UA-Platform")),
		Mobile:          stringPtr(headerValue(r, "Sec-CH-UA-Mobile")),
		Model:           stringPtr(headerValue(r, "Sec-CH-UA-Model")),
		Arch:            stringPtr(headerValue(r, "Sec-CH-UA-Arch")),
		Bitness:         stringPtr(headerValue(r, "Sec-CH-UA-Bitness")),
		FullVersionList: stringPtr(headerValue(r, "Sec-CH-UA-Full-Version-List")),
		PlatformVersion: stringPtr(headerValue(r, "Sec-CH-UA-Platform-Version")),
	}

	if !hasPtr(info.UA, info.Platform, info.Mobile, info.Model, info.Arch, info.Bitness, info.FullVersionList, info.PlatformVersion) {
		return nil
	}

	return &info
}

func collectRequestHeaders(r *http.Request) []HeaderEntry {
	if r == nil {
		return nil
	}

	entries := make([]HeaderEntry, 0, len(r.Header)+2)

	if host := strings.TrimSpace(r.Host); host != "" {
		entries = append(entries, HeaderEntry{Key: "Host", Value: host})
	}

	if r.ContentLength > 0 {
		entries = append(entries, HeaderEntry{Key: "Content-Length", Value: strconv.FormatInt(r.ContentLength, 10)})
	}

	if len(r.TransferEncoding) > 0 {
		entries = append(entries, HeaderEntry{Key: "Transfer-Encoding", Value: strings.Join(r.TransferEncoding, ", ")})
	}

	for key, values := range r.Header {
		if len(values) == 0 {
			continue
		}

		cleaned := make([]string, 0, len(values))
		for _, value := range values {
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}
			cleaned = append(cleaned, value)
		}

		if len(cleaned) == 0 {
			continue
		}

		entries = append(entries, HeaderEntry{Key: key, Value: strings.Join(cleaned, ", ")})
	}

	if len(entries) == 0 {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Key) < strings.ToLower(entries[j].Key)
	})

	return entries
}

func detectScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}

	if proto := firstHeaderToken(headerValue(r, "X-Forwarded-Proto")); proto != "" {
		return proto
	}

	if proto := forwardedParam(headerValue(r, "Forwarded"), "proto"); proto != "" {
		return proto
	}

	if r.URL != nil && r.URL.Scheme != "" {
		return r.URL.Scheme
	}

	if r.Proto != "" {
		return "http"
	}

	return ""
}

func headerValue(r *http.Request, key string) string {
	values := r.Header.Values(key)
	if len(values) == 0 {
		return ""
	}

	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		cleaned = append(cleaned, value)
	}

	if len(cleaned) == 0 {
		return ""
	}

	return strings.Join(cleaned, ", ")
}

func firstHeaderToken(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}

func forwardedParam(value, key string) string {
	if value == "" {
		return ""
	}

	entries := strings.Split(value, ",")
	for _, entry := range entries {
		params := strings.Split(entry, ";")
		for _, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				continue
			}

			parts := strings.SplitN(param, "=", 2)
			if len(parts) != 2 {
				continue
			}

			if !strings.EqualFold(strings.TrimSpace(parts[0]), key) {
				continue
			}

			return strings.Trim(strings.TrimSpace(parts[1]), "\"")
		}
	}

	return ""
}

func hasPtr(values ...*string) bool {
	for _, value := range values {
		if value != nil {
			return true
		}
	}
	return false
}

func stringPtr(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
