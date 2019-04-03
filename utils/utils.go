package utils

import (
	"net/http"
	"strings"
)

// GetSubdomain : return subdomain of request
func GetSubdomain(r *http.Request) string {
	fullHost := r.Host
	splitHost := strings.Split(fullHost, ".")
	if len(splitHost) < 1 {
		return ""
	}
	return splitHost[0]
}
