package middlewares

import (
	//"golang.org/x/exp/slices"

	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var allowedTypes = []string{
	"application/javascript",
	"application/json",
	"application/gzip",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandle compresses data with gzip
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			body, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = body
		} else {
			next.ServeHTTP(w, r)
			return
		}

		gzw, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer gzw.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gzw}, r)
	})
}
