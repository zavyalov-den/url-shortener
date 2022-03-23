package middlewares

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
)

func CompressMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: gzip
		w.Header()

		next.ServeHTTP(w, r)
	})
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to buffer: %s", err)
	}

	err = w.Close()

	return b.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to init reader: %s", err)
	}
	defer r.Close()

	var b bytes.Buffer

	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read from buffer: %s", err)
	}

	return b.Bytes(), nil
}
