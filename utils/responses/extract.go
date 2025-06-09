package responses

import (
	"net"
	"net/http"
	"strings"
)

// getClientIP attempts to get the real client IP address from HTTP headers or RemoteAddr.
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (may contain multiple IPs)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		// Take the first valid IP address
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		ip := strings.TrimSpace(xRealIP)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Fallback: parse IP from RemoteAddr (host:port)
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Last fallback: return RemoteAddr as-is (may include port)
	return r.RemoteAddr
}

// extractRequestInfo extracts relevant request information as a struct.
func extractRequestInfo(r *http.Request) RequestInfo {
	return RequestInfo{
		Method:    r.Method,
		Path:      r.URL.Path,
		UserAgent: r.UserAgent(),
		RemoteIP:  getClientIP(r),
	}
}
