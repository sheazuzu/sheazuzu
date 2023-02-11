/*
 * debug_test.go
 * Created on 25.02.2020
 * Copyright (C) 2020 Volkswagen AG, All rights reserved
 *
 */

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebugHandler(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, DebugHandler())
}
