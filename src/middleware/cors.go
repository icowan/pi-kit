/**
 * @Time : 2022/10/19 3:40 PM
 * @Author : solacowa@gmail.com
 * @File : cors
 * @Software: GoLand
 */

package middleware

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"net/http"
)

func CORSMethodMiddleware(logger log.Logger, corsHeaders map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, val := range corsHeaders {
				w.Header().Set(key, val)
			}
			w.Header().Set("Connection", "keep-alive")
			if r.Method == "OPTIONS" {
				return
			}
			_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
			next.ServeHTTP(w, r)
		})
	}
}
