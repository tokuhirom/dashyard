package acl

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IPAllowList holds parsed CIDR networks that are allowed access.
type IPAllowList struct {
	networks []*net.IPNet
}

// New parses the given IP/CIDR strings and returns an IPAllowList.
// Returns nil if the allow list is empty (meaning allow all).
// Returns an error if any entry is invalid.
func New(allow []string) (*IPAllowList, error) {
	if len(allow) == 0 {
		return nil, nil
	}

	networks := make([]*net.IPNet, 0, len(allow))
	for _, entry := range allow {
		_, ipNet, err := net.ParseCIDR(entry)
		if err != nil {
			// Try parsing as a single IP and normalize to CIDR
			ip := net.ParseIP(entry)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP/CIDR %q: %w", entry, err)
			}
			if ip.To4() != nil {
				_, ipNet, _ = net.ParseCIDR(ip.String() + "/32")
			} else {
				_, ipNet, _ = net.ParseCIDR(ip.String() + "/128")
			}
		}
		networks = append(networks, ipNet)
	}

	return &IPAllowList{networks: networks}, nil
}

// Contains reports whether the given IP string is in the allow list.
func (a *IPAllowList) Contains(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	for _, n := range a.networks {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// Middleware returns a Gin middleware that rejects requests from IPs not in the
// allow list. If allowList is nil, all requests are allowed.
func Middleware(allowList *IPAllowList) gin.HandlerFunc {
	return func(c *gin.Context) {
		if allowList == nil {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		if !allowList.Contains(clientIP) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
