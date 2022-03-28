package handler

import (
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"net/http"
	//"github.com/jackc/pgx/v4"
)

func Ping(db *storage.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		//conn, err := pgx.Connect

		w.WriteHeader(http.StatusOK)
	}
}
