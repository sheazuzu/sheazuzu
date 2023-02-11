/*
 * logger_test.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
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

func TestLoggerHandler(t *testing.T) {
	t.Parallel()

	logger := GetLogger("error", "json")
	var handler http.Handler
	handler = Logger(logger)(http.NotFoundHandler())
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	assert.NotNil(t, handler)
}
