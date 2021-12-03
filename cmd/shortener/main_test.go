package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Want struct {
	statusCode int
	body       string
}

func Test_handler(t *testing.T) {
	tests := []struct {
		name   string
		urls   map[string]string
		method string
		body   string
		params string
		want   Want
	}{
		{
			"shorten",
			getURLs(false),
			http.MethodPost,
			"https://yandex.ru",
			"",
			Want{
				statusCode: 201,
				body:       "http://localhost:8080/e9db20b2",
			},
		},
		{
			"expand",
			getURLs(true),
			http.MethodGet,
			"",
			"/e9db20b2",
			Want{
				statusCode: 307,
				body:       "",
			},
		},
		{
			"returns 404 on url that doesn't exist",
			getURLs(false),
			http.MethodGet,
			"",
			"/asdfa",
			Want{
				statusCode: 404,
				body:       "404 page not found\n",
			},
		},
		{
			"returns 404 on invalid request URL",
			getURLs(false),
			http.MethodGet,
			"",
			"/wrong/url",
			Want{
				statusCode: 404,
				body:       "404 page not found\n",
				//body:       "invalid request URL\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := newTestServer(tt.urls)

			cl := ts.Client()
			cl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			var resp *http.Response

			switch tt.method {
			case http.MethodGet:
				resp, _ = cl.Get(ts.URL + tt.params)
			case http.MethodPost:
				resp, _ = cl.Post(ts.URL+tt.params, "text/plain; charset=utf8", strings.NewReader(tt.body))
			default:
				t.Fatal("Method is not allowed")
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				assert.NoError(t, err, "can't read response body")
			}

			defer resp.Body.Close()

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

func newTestServer(urls map[string]string) *httptest.Server {
	r := chi.NewRouter()

	r.Get("/{shortUrl}", getHandler(urls))
	r.Post("/", postHandler(urls))

	return httptest.NewServer(r)
}
