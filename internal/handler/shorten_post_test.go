package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zavyalov-den/url-shortener/internal/config"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_shortenPost(t *testing.T) {
	tests := []struct {
		name string
		db   *storage.DB
		body string
		want want
	}{
		{
			"shorten",
			newTestDB(false),
			`{"url": "https://yandex.ru"}`,
			want{
				statusCode: 201,
				body:       fmt.Sprintf(`{"result":"%s/e9db20b2"}`, config.Conf.BaseURL),
			},
		},
		{
			"shorten negative",
			newTestDB(false),
			"",
			want{
				statusCode: 400,
				body:       "{\"error\":\"request body is not a valid json\"}\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newShortenPostTestServer(tt.db)

			cl := ts.Client()

			var resp *http.Response

			resp, err := cl.Post(ts.URL+"/api/shorten", "application/json", strings.NewReader(tt.body))
			require.NoError(t, err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "can't read response body")

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Wrong status code")
			assert.Equal(t, tt.want.body, string(respBody), "Wrong status code")

		})
	}
}

func newShortenPostTestServer(db *storage.DB) *httptest.Server {
	r := chi.NewRouter()

	r.Post("/api/shorten", ShortenPost(db))

	return httptest.NewServer(r)
}
