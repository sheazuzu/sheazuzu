/*
 *  logger.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

// Package logging provides constants and functions for a consistent logger configuration within the NGW application landscape.
package logging

import (
	"net/http"
	"os"
	"runtime/debug"
	"sheazuzu/common/src/tracing"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger returns a new logger. The parameter level defines which messages will be logged. Possible are "info", "debug", "warn" or "error".
// The parameter format is used to select an encoding and build the encoder config. Valid values for format are "minimal", "json" or "console".
// A format = minimal e.g. prevents stacktraces from being logged.
func GetLogger(level, format string) *zap.SugaredLogger {

	loggerCfg := zap.NewProductionConfig()

	switch strings.ToLower(level) {
	case "info":
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "debug":
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	if format == "minimal" {
		loggerCfg.Encoding = "console"
	} else {
		loggerCfg.Encoding = format
	}

	switch format {
	case "console":
		loggerCfg.EncoderConfig = zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "@timestamp",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		}
	case "json":
		loggerCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		loggerCfg.EncoderConfig.TimeKey = "@timestamp" // NGWDEV-257: Unify the timestamp attribute name
		loggerCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	case "minimal":
		loggerCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		loggerCfg.EncoderConfig.TimeKey = ""
		loggerCfg.EncoderConfig.CallerKey = ""
		loggerCfg.DisableCaller = true
		loggerCfg.DisableStacktrace = true
	}

	// cant return an error, because no user input config is permitted that isnt checked beforehand
	logger, _ := loggerCfg.Build()

	return logger.Sugar()
}

// RequestLogHandler returns a middleware that logs incoming http calls and their duration
// It also adds the traceId to the logging context, but make sure the `tracing.TraceHandler` middleware is in the chain before this handler
func RequestLogHandler(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := WithLogger(r.Context(), logger)
			defer RecoverPanic(logger)

			// add the traceid to the logger
			logger := tracing.LoggerWithTraceID(ctx, logger)
			r = r.WithContext(ctx)

			start := time.Now()

			inner.ServeHTTP(w, r)

			logger.Debugf(
				"%s %s %s",
				r.Method,
				r.RequestURI,
				time.Since(start),
			)
		})
	}
}

var exitFn = os.Exit

// RecoverPanic logs the panic which occurred in a goroutine and exits the routine.
func RecoverPanic(logger *zap.SugaredLogger) {
	if r := recover(); r != nil {
		logger.Errorf("Panic: %v,\n%s", r, debug.Stack())
		exitFn(1)
	}
}
