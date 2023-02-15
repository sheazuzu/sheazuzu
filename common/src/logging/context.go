/*
 *  context.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package logging

import (
	"context"
	"go.uber.org/zap"
)

type logCtxKey struct{}

// WithLogger adds the logger to the context.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, logCtxKey{}, logger)
}

// ContextLogger retrieves the logger from the context.
// If the context does not contain a logger, a new no-op logger will be returned.
func ContextLogger(ctx context.Context) *zap.SugaredLogger {
	logger, _ := ctx.Value(logCtxKey{}).(*zap.SugaredLogger)
	if logger == nil {
		return zap.NewNop().Sugar() // make sure no nil pointer occur
	}
	return logger
}

// ContextLoggerWith adds the arguments to the logger included in the context. Afterwards, it returns the context with the updated logger and the logger itself.
// If the context does not contain a logger, a new no-op logger (without the arguments set) will be returned besides the context.
func ContextLoggerWith(ctx context.Context, args ...interface{}) (context.Context, *zap.SugaredLogger) {
	logger, _ := ctx.Value(logCtxKey{}).(*zap.SugaredLogger)
	if logger == nil {
		return ctx, zap.NewNop().Sugar() // make sure no nil pointer occur
	}
	logger = logger.With(args...)
	return WithLogger(ctx, logger), logger
}
