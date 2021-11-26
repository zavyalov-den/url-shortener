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

var urls = make(map[string]string)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			break
		}
		url := string(data)

		md5 := crypto.MD5.New()
		md5.Write(data)
		short := hex.EncodeToString(md5.Sum(nil))[:8]

		urls[short] = url

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(short))
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			break
		}
	case http.MethodGet:
		pattern := "^/?[0-9a-f]+/?$"
		re := regexp.MustCompile(pattern)
		url := r.URL.String()
		match := re.Match([]byte(url))
		if !match {
			fmt.Println(match, url)
			http.Error(w, "invalid request", http.StatusBadRequest)
			break
		}

		shortURL := strings.Replace(re.FindStringSubmatch(url)[0], "/", "", -1)
		longURL := urls[shortURL]
		if longURL == "" {
			//http.NotFound(w, r)
			http.Error(w, "invalid request", http.StatusBadRequest)
			break
		}

		//w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//_, err := w.Write(nil)
		//if err != nil {
		//	http.Error(w, "invalid request", http.StatusBadRequest)
		//	break
		//}

	default:
		http.Error(w, "invalid request", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
