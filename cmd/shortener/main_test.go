package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPage(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		want    want
		request string
		longurl string
	}{
		{
			name: "test #1",
			want: want{
				statusCode:  201,
				response:    "http://localhost/123",
				contentType: "text/plain",
			},
			request: "/",
			longurl: "321",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, test.request, nil)
			w := httptest.NewRecorder()
			MainPage(w, request)
			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)

		})
	}
}

func TestSubPage(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "test #2",
			want: want{
				statusCode: 307,
			},
			request: "/123",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, test.request, nil)
			w := httptest.NewRecorder()
			SubPage(w, request)
			res := w.Result()

			assert.Equal(t, test.want.statusCode, res.StatusCode)
		})
	}
}
