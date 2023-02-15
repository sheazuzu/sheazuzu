/*
 * trace_test.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package tracing

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func TestRegisterExporter(t *testing.T) {

	type input struct {
		serviceName string
		agent       string
		collector   string
	}

	type output struct {
		errorExpected bool
		errorText     string
	}

	type test struct {
		input  input
		output output
	}

	cases := map[string]test{
		"everything set": {
			input: input{
				serviceName: "test-service",
				agent:       "localhost:6789",
				collector:   "http://localhost:1234/api/traces",
			},
			output: output{
				errorExpected: false,
			},
		},
		"nothing set": {
			input: input{},
			output: output{
				errorExpected: true,
				errorText:     "missing endpoint for Jaeger exporter",
			},
		},
		"agent only": {
			input: input{
				agent: "localhost:6831",
			},
			output: output{
				errorExpected: false,
			},
		},
		"collector only": {
			input: input{
				collector: "localhost:6831/api/traces",
			},
			output: output{
				errorExpected: false,
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert := assert.New(t)

			exporter, err := RegisterTraceExporter(tc.input.agent, tc.input.collector, tc.input.serviceName)
			if tc.output.errorExpected {
				assert.EqualError(err, tc.output.errorText, "wrong error")
			} else {
				assert.NotNil(exporter)
				assert.Nil(err, "expected error to be nil, but it wasn`t")
			}
		})
	}
}

func TestStartSpan(t *testing.T) {
	t.Parallel()

	ctx, span := StartSpan(context.Background(), "test")

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
}

func TestTraceHandler(t *testing.T) {
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
		"204": {
			input: input{
				status: 204,
			},
		},
		"400": {
			input: input{
				status: 400,
			},
		},
		"401": {
			input: input{
				status: 401,
			},
		},
		"404": {
			input: input{
				status: 404,
			},
		},
		"429": {
			input: input{
				status: 429,
			},
		},
		"500": {
			input: input{
				status: 500,
			},
		},
		"123": {
			input: input{
				status: 123,
			},
		},
	}

	getHandler := func(status int) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write([]byte("test"))
		}
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			finalHandler := TraceHandler(zap.NewNop().Sugar(), "test")(http.HandlerFunc(getHandler(tc.input.status)))

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://localhost:8080/irrelevantUrl", ioutil.NopCloser(bytes.NewBufferString("test")))

			finalHandler.ServeHTTP(recorder, request)

			assert.Equal(t, "test", recorder.Body.String())
		})
	}
}

func TestAddStatus(t *testing.T) {
	t.Parallel()

	type input struct {
		err    error
		status int
	}

	type test struct {
		input input
	}

	cases := map[string]test{
		"with error": {
			input: input{
				err:    fmt.Errorf("mocked error"),
				status: 500,
			},
		},
		"without error": {
			input: input{
				err:    nil,
				status: 200,
			},
		},
	}

	for name, tc := range cases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			span := &trace.Span{}

			AddStatus(span, tc.input.err, tc.input.status)
		})
	}
}
