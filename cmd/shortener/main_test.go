package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var positiveUsers = map[string]string{
	"e9db20b2": "https://yandex.ru",
}

type Want struct {
	statusCode int
	body       string
}

func Test_handler(t *testing.T) {
	tests := []struct {
		name        string
		requestURL  string
		requestBody string
		method      string
		users       map[string]string
		want        Want
	}{
		{
			"shorten",
			"/",
			"https://yandex.ru",
			http.MethodPost,
			make(map[string]string),
			Want{
				statusCode: 201,
				body:       "http://localhost:8080/e9db20b2",
			},
		},
		{
			"expand",
			"/e9db20b2",
			"",
			http.MethodGet,
			positiveUsers,
			Want{
				statusCode: 307,
				body:       "",
			},
		},
		{
			"returns 404 on url that doesn't exist",
			"/e9db20b2",
			"",
			http.MethodGet,
			make(map[string]string),
			Want{
				statusCode: 404,
				body:       "404 page not found\n",
			},
		},
		{
			"returns 400 on invalid request URL",
			"/e9db20b2/asdf/adf",
			"",
			http.MethodGet,
			make(map[string]string),
			Want{
				statusCode: 400,
				body:       "invalid request URL\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test
			request := httptest.NewRequest(tt.method, tt.requestURL, strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			h := handler(tt.users)

			h.ServeHTTP(w, request)

			resp := w.Result()

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
