package service

import (
	"crypto"
	"encoding/hex"
	"github.com/zavyalov-den/url-shortener/internal/config"
)

func Shorten(data []byte) string {
	md5 := crypto.MD5.New()
	md5.Write(data)
	short := hex.EncodeToString(md5.Sum(nil))[:8]

	return short
}

func ShortToURL(s string) string {
	return config.GetConfigInstance().BaseURL + "/" + s
}
