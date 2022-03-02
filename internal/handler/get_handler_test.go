package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name   string
		urls   map[string]string
		body   string
		params string
		want   want
	}{
		{
			"expand",
			getURLs(true),
			"",
			"/e9db20b2",
			want{
				statusCode: 307,
				body:       "",
			},
		},
		{
			"returns 404 on url that doesn't exist",
			getURLs(false),
			"",
			"/asdfa",
			want{
				statusCode: 404,
				body:       "404 page not found\n",
			},
		},
		{
			"returns 404 on invalid request URL",
			getURLs(false),
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
			ts := newGetTestServer(tt.urls)

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

func getURLs(notEmpty bool) map[string]string {
	if notEmpty {
		return map[string]string{
			"e9db20b2": "https://yandex.ru",
		}
	}
	return make(map[string]string)
}

func newGetTestServer(urls map[string]string) *httptest.Server {
	r := chi.NewRouter()

	r.Get("/{shortUrl}", Get(urls))

	return httptest.NewServer(r)
}
