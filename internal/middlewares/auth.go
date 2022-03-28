package middlewares

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"net/http"
	"strconv"
)

var (
	currentUserId = 0
	nonce         []byte
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for cookie, if exists set UserID as context?
		// if cookie is absent - generate a new one.

		cookie, err := r.Cookie("auth")
		if errors.Is(err, http.ErrNoCookie) {
			cookie = createAuthCookie()
		} else if err != nil {
			cookie = createAuthCookie()
			//next.ServeHTTP(w, r)
			//return
		}

		userID := decodeAuthCookie(cookie)
		if userID == 0 {
			cookie = createAuthCookie()
			//next.ServeHTTP(w, r)
			//return
		}
		ctx := context.WithValue(nil, "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

func decodeAuthCookie(cookie *http.Cookie) int {
	key := sha256.Sum256([]byte(config.C.AuthKey))

	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Println("new cipher", err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		fmt.Println("new gcm", err)
	}

	nonce := getNonce(aesGCM.NonceSize())
	src, err := aesGCM.Open(nil, nonce, []byte(cookie.Value), nil)
	if err != nil {
		fmt.Println("open", err)
	}

	fmt.Println(string(src))
	userID, err := strconv.Atoi(string(src))
	if err != nil {
		return 0
	}
	return userID
}

func createAuthCookie() *http.Cookie {
	currentUserId++

	key := sha256.Sum256([]byte(config.C.AuthKey))

	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Println("new cipher", err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		fmt.Println("new gcm", err)
	}

	nonce := getNonce(aesGCM.NonceSize())

	sealedCookie := aesGCM.Seal(nil, nonce, []byte(strconv.Itoa(currentUserId)), nil)

	return &http.Cookie{
		Name:  "auth",
		Value: string(sealedCookie),
	}
}

func getNonce(n int) []byte {
	var b []byte
	if !bytes.Equal(nonce, b) {
		return nonce
	}

	nonce = make([]byte, n)
	_, err := rand.Read(nonce)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nonce
}
