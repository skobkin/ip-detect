package clientinfo

import (
	"context"
	"net"
	"strings"
	"time"
)

const fallbackLookupTimeout = 200 * time.Millisecond

func reverseLookup(ctx context.Context, ip string, timeout time.Duration) string {
	if timeout <= 0 {
		timeout = fallbackLookupTimeout
	}

	resolverCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	names, err := net.DefaultResolver.LookupAddr(resolverCtx, ip)
	if err != nil || len(names) == 0 {
		return ""
	}

	return strings.TrimSuffix(names[0], ".")
}
