package main

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func handler(urls map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			data, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "invalid requestURL", http.StatusBadRequest)
				break
			}
			url := string(data)

			md5 := crypto.MD5.New()
			md5.Write(data)
			short := hex.EncodeToString(md5.Sum(nil))[:8]

			urls[short] = url

			w.WriteHeader(http.StatusCreated)
			_, err = w.Write([]byte("http://localhost:8080/" + short))
			if err != nil {
				http.Error(w, "invalid requestURL", http.StatusBadRequest)
				break
			}
		case http.MethodGet:
			pattern := "^/?[0-9a-f]+/?$"
			re := regexp.MustCompile(pattern)
			url := r.URL.String()
			match := re.Match([]byte(url))
			if !match {
				fmt.Println(match, url)
				http.Error(w, "invalid request URL", http.StatusBadRequest)
				break
			}

			shortURL := strings.Replace(re.FindStringSubmatch(url)[0], "/", "", -1)
			longURL := urls[shortURL]
			if longURL == "" {
				http.NotFound(w, r)
				//http.Error(w, "invalid requestURL", http.StatusBadRequest)
				break
			}

			w.Header().Set("Location", longURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
		default:
			http.Error(w, "invalid requestURL", http.StatusBadRequest)
		}
	}

}

func main() {
	urls := make(map[string]string)

	http.HandleFunc("/", handler(urls))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
