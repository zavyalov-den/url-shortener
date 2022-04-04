package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zavyalov-den/url-shortener/internal/service"
	"github.com/zavyalov-den/url-shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type want struct {
	statusCode int
	body       string
}

func Test_GetHandler(t *testing.T) {
	tests := []struct {
		name       string
		db         storage.Storage
		dbTestData bool
		body       string
		params     string
		want       want
	}{
		{
			"expand",
			newTestDB(),
			false,
			"",
			"/e9db20b2",
			want{
				statusCode: 307,
				body:       "",
			},
		},
		{
			"returns 404 on url that doesn't exist",
			newTestDB(),
			true,
			"",
			"/asdfa",
			want{
				statusCode: 404,
				body:       "404 page not found\n",
			},
		},
		{
			"returns 404 on invalid request URL",
			newTestDB(),
			true,
			"",
			"/wrong/url",
			want{
				statusCode: 404,
				body:       "404 page not found\n",
				//body:       "invalid request URL\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbTestData {
				err := tt.db.SaveURL(1, storage.UserURL{
					ShortURL:    service.ShortToURL("e9db20b2"),
					OriginalURL: "",
				})
				assert.NoError(t, err)
			}

			ts := newGetTestServer(tt.db)

			cl := ts.Client()
			cl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			resp, err := cl.Get(ts.URL + tt.params)
			require.NoError(t, err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				assert.NoError(t, err, "can't read response body")
			}

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Wrong status code")
			assert.Equal(t, tt.want.body, string(respBody), "Wrong status code")
		})
	}
}

func newTestDB() storage.Storage {
	db := storage.NewStorage()
	return db
}

func newGetTestServer(db storage.Storage) *httptest.Server {
	r := chi.NewRouter()

	r.Get("/{shortUrl}", GetFullURL(db))

	return httptest.NewServer(r)
}
