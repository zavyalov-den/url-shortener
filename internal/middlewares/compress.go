package middlewares

import "net/http"

func CompressMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: gzip

		next.ServeHTTP(w, r)
	})
}
