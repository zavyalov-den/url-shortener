package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_PostHandler(t *testing.T) {
	tests := []struct {
		name   string
		urls   map[string]string
		body   string
		params string
		want   want
	}{
		{
			"shorten",
			getURLs(false),
			"https://yandex.ru",
			"",
			want{
				statusCode: 201,
				body:       "http://localhost:8080/e9db20b2",
			},
		},
		{
			"shorten negative",
			getURLs(false),
			"",
			"",
			want{
				statusCode: 400,
				body:       "request body must not be empty\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newPostTestServer(tt.urls)

			cl := ts.Client()
			cl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			var resp *http.Response

			resp, err := cl.Post(ts.URL+tt.params, "text/plain; charset=utf8", strings.NewReader(tt.body))
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

func newPostTestServer(urls map[string]string) *httptest.Server {
	r := chi.NewRouter()

	r.Post("/", Post(urls))

	return httptest.NewServer(r)
}
