/*
 *  logger_test.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	t.Parallel()

	type test struct {
		format      string
		level       string
		nilExpected bool
	}

	tests := map[string]test{
		"log level debug": {
			format:      "json",
			level:       "debug",
			nilExpected: false,
		},
		"log level info": {
			format:      "json",
			level:       "info",
			nilExpected: false,
		},
		"log level warn": {
			format:      "json",
			level:       "warn",
			nilExpected: false,
		},
		"log level error": {
			format:      "json",
			level:       "error",
			nilExpected: false,
		},
		"log format console": {
			format:      "console",
			level:       "error",
			nilExpected: false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger := GetLogger(tc.level, tc.format)
			if tc.nilExpected && logger != nil {
				t.Fatal("expected the logger to be nil, but it wasn`t")
			}

			if !tc.nilExpected && logger == nil {
				t.Fatal("expected the logger not to be nil, but it was")
			}
		})
	}
}

func TestRequestLogHandler(t *testing.T) {
	t.Parallel()

	logger := GetLogger("error", "json")
	var handler http.Handler
	handler = RequestLogHandler(logger)(http.NotFoundHandler())
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	assert.NotNil(t, handler)
}

func TestRecoverPanic(t *testing.T) {

	logger := GetLogger("debug", "minimal")

	defer RecoverPanic(logger)

	exitCalled := false

	exitFn = func(code int) {
		assert.Equal(t, 1, code)
		exitCalled = true
	}

	panic("Don't Panic")

	//goland:noinspection GoUnreachableCode
	assert.True(t, exitCalled)
}
