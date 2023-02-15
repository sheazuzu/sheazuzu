/*
 *  context_test.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package logging

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

// this test shows how logger can be nested
func TestContextLogger(t *testing.T) {

	assert.NotNil(t, ContextLogger(context.Background()))
	ctx, l := ContextLoggerWith(context.Background())
	assert.NotNil(t, ctx)
	assert.NotNil(t, l)

	rootCtx := WithLogger(context.Background(), GetLogger("info", "console"))
	assert.NotNil(t, rootCtx)

	rootLogger := ContextLogger(rootCtx)
	assert.NotNil(t, rootLogger)
	rootLogger.Info("Root Logger")

	childCtx1 := WithLogger(rootCtx, rootLogger.With("child1", "logger1"))
	assert.NotNil(t, childCtx1)

	childLogger1 := ContextLogger(childCtx1)
	assert.NotNil(t, childLogger1)

	childLogger1.Info("Child Logger 1")

	childCtx2, childLogger2 := ContextLoggerWith(childCtx1, "child2", "logger2")
	assert.NotNil(t, childCtx2)
	assert.NotNil(t, childLogger2)

	childLogger2.Info("Child Logger 2")

	rootLogger.Info("Root Logger")
	childLogger1.Info("Child Logger 1")
}

func ExampleWithLogger() {
	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	fmt.Println(WithLogger(ctx, logger))
	// Output: context.Background.WithValue(type logging.logCtxKey, val <not Stringer>)
}

func ExampleContextLogger() {
	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	loggerFromContext := ContextLogger(WithLogger(ctx, logger))
	fmt.Println(reflect.DeepEqual(logger, loggerFromContext))
	// Output: true
}

func ExampleContextLoggerWith() {
	logger := zap.NewNop().Sugar()
	ctx := WithLogger(context.Background(), logger)
	ContextLoggerWith(ctx, "test1", "test2")
}
