/*
 *  ui_test.go
 *  Created on 29.11.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package swagger

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
)

func TestRegisterSwaggerHandlers(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	RegisterSwaggerHandlers(r, &openapi3.Swagger{}, "myContextPath")

}
