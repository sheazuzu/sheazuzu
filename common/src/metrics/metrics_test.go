/*
 * metrics_test.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package metrics

import (
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	t.Parallel()

	err := RegisterHandler("test_namespace", chi.NewRouter())

	assert.Nil(t, err)
}

func TestGetMetricsRecordingHandler(t *testing.T) {

	t.Parallel()

	type input struct {
		status int
	}

	type test struct {
		input input
	}

	cases := map[string]test{
		"200": {
			input: input{
				status: 200,
			},
		},
		"300": {
			input: input{
				status: 300,
			},
		},
		"400": {
			input: input{
				status: 400,
			},
		},
		"500": {
			input: input{
				status: 500,
			},
		},
		"600": {
			input: input{
				status: 600,
			},
		},
	}

	getHandler := func(status int) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write([]byte("test"))
		})
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://localhost:8080/irrelevantUrl", nil)

			handler := GetMetricsRecordingHandler(getHandler(tc.input.status))

			handler.ServeHTTP(recorder, request)
		})
	}
}
