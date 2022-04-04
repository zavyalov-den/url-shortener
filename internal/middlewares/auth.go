package middlewares

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"net/http"
	"strconv"
	"time"
)

type CryptoSvc struct {
	aesBlock   cipher.Block
	aesGCM     cipher.AEAD
	nonce      []byte
	lastUserID int
}

var (
	cryptoSvc *CryptoSvc
	n         int
)

func GetCryptoSvcInstance() *CryptoSvc {
	if cryptoSvc != nil {
		return cryptoSvc
	}

	n++
	fmt.Println("creating crypto instance #", n)

	key := sha256.Sum256([]byte(config.GetConfigInstance().AuthKey))

	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		//fmt.Println("new cipher", err)
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(aesBlock)
	if err != nil {
		//fmt.Println("new gcm", err)
		panic(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	//_, err = rand.Read(nonce)
	//if err != nil {
	//	fmt.Println("read err: ", err.Error())
	//	fmt.Println(err.Error())
	//}

	cryptoSvc = &CryptoSvc{
		aesBlock:   aesBlock,
		aesGCM:     aesGCM,
		nonce:      nonce,
		lastUserID: 0,
	}
	return cryptoSvc
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := GetCryptoSvcInstance()

		cookie, err := r.Cookie("auth")
		if err != nil {
			cookie = c.createAuthCookie()
		}
		//if err = cookie.Valid(); err != nil {
		//	fmt.Println("cookie is not valid: ", err)
		//	cookie = c.createAuthCookie()
		//}

		userID := c.decodeAuthCookie(cookie)
		if userID == 0 {
			cookie = c.createAuthCookie()
			userID = c.decodeAuthCookie(cookie)
		}
		ctx := context.WithValue(r.Context(), "auth", userID)

		http.SetCookie(w, cookie)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *CryptoSvc) decodeAuthCookie(cookie *http.Cookie) int {
	cookieBytes, err := hex.DecodeString(cookie.Value)
	if err != nil {
		fmt.Println("failed to decode a cookie :(", err)
		return 0
	}

	src, err := c.aesGCM.Open(nil, c.nonce, cookieBytes, nil)
	if err != nil {
		fmt.Println("gcm open failed: ", err)
		return 0
	}

	userID, err := strconv.Atoi(string(src))
	if err != nil {
		return 0
	}
	return userID
}

func (c *CryptoSvc) createAuthCookie() *http.Cookie {
	c.lastUserID++

	byteString := hex.EncodeToString([]byte(strconv.Itoa(c.lastUserID)))

	sealedCookie := c.aesGCM.Seal(nil, c.nonce, []byte(byteString), nil)

	return &http.Cookie{
		Name:    "auth",
		Value:   hex.EncodeToString(sealedCookie),
		Expires: time.Now().Add(8 * time.Hour),
		Path:    "/",
	}
}
