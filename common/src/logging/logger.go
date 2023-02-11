/*
 * logger.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package logging

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger returns a new logger
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

// Logger returns a middleware that logs incoming http calls and their duration
func Logger(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
