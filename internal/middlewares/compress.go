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
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r.Body = gzr
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
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
		//if !strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip") {
		//	next.ServeHTTP(w, r)
		//	return
		//}
		//
		////// decode body
		////body, err := gzip.NewReader(r.Body)
		////if err != nil {
		////	http.Error(w, err.Error(), http.StatusBadRequest)
		////	return
		////}
		////r.Body = body
		//
		//gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//	return
		//}
		//defer gz.Close()
		//
		//w.Header().Set("Content-Encoding", "gzip")
		//
		//next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		//contentType := strings.ToLower(r.Header.Get("Content-Type"))
		//
		//var allowedCT bool
		//
		//// would love to replace it with golang.org/x/exp/slices.Contains()
		//// but tests run on 1.17 :(
		//for _, ct := range allowedTypes {
		//	if ct == contentType {
		//		allowedCT = true
		//	}
		//}
		//
		//if !allowedCT {
		//	next.ServeHTTP(w, r)
		//	return
		//}
		//
	})
}
