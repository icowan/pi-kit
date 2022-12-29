/**
 * @Time : 2022/10/19 3:35 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package middleware

import (
	"fmt"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		fmt.Println("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		//_ = level.Info(logger).Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
