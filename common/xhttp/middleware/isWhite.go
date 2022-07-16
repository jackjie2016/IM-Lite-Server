package middleware

import (
	"net/http"
	"strings"
)

func (m *AuthMiddleware) isWhite(req *http.Request) bool {
	return strings.Contains(req.URL.Path, "/white/")
}
