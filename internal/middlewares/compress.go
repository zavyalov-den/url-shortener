package middlewares

import (
	"compress/gzip"
	//"golang.org/x/exp/slices"

	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedTypes := []string{
			"application/javascript",
			"application/json",
			"text/css",
			"text/html",
			"text/plain",
			"text/xml",
		}

		//!slices.Contains(allowedTypes, w.Header().Get("Content-Type")) {
		if !strings.Contains(strings.ToLower(w.Header().Get("Accept-Encoding")), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		contentType := w.Header().Get("Content-Type")

		var allowedCT bool

		// would love to replace it with golang.org/x/exp/slices.Contains() but tests run on 1.17 :(
		for _, ct := range allowedTypes {
			if ct == contentType {
				allowedCT = true
			}
		}

		if !allowedCT {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(&gzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}, r)
	})
}

//func Compress(data []byte) ([]byte, error) {
//	var b bytes.Buffer
//
//	w := gzip.NewWriter(&b)
//
//	_, err := w.Write(data)
//	if err != nil {
//		return nil, fmt.Errorf("failed to write data to buffer: %s", err)
//	}
//
//	err = w.Close()
//
//	return b.Bytes(), nil
//}
//
//func Decompress(data []byte) ([]byte, error) {
//	r, err := gzip.NewReader(bytes.NewReader(data))
//	if err != nil {
//		return nil, fmt.Errorf("failed to init reader: %s", err)
//	}
//	defer r.Close()
//
//	var b bytes.Buffer
//
//	_, err = b.ReadFrom(r)
//	if err != nil {
//		return nil, fmt.Errorf("failed to read from buffer: %s", err)
//	}
//
//	return b.Bytes(), nil
//}
