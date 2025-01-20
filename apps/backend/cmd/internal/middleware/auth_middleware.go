package middleware

import "net/http"

type AuthMiddleware struct {
	Handler http.Handler
}
