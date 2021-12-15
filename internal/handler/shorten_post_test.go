package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_shortenPost(t *testing.T) {
	// todo

	tests := []struct {
		name string
		urls map[string]string
		body string
		want want
	}{
		{
			"shorten",
			getURLs(false),
			`{"url": "https://yandex.ru"}`,
			want{
				statusCode: 201,
				body:       `{"result": "http://localhost:8080/e9db20b2"}`,
			},
		},
		{
			"shorten negative",
			getURLs(false),
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

			var resp *http.Response

			resp, err := cl.Post(ts.URL, "application/json", strings.NewReader(tt.body))
			require.NoError(t, err)

			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "can't read response body")

			assert.Equal(t, tt.want.statusCode, resp.StatusCode, "Wrong status code")
			assert.Equal(t, tt.want.body, string(respBody), "Wrong status code")

		})
	}
}
