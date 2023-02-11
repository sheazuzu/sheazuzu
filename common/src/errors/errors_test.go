/*
 *  errors_test.go
 *  Created on 29.11.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testErrorInfoProviderPointerReceiver struct {
}

func (t *testErrorInfoProviderPointerReceiver) ErrorInfo() Info {
	return Info{
		Name: "outer-info",
		Val: []Info{
			{
				Name: "inner-info",
				Val:  "inner-val",
			},
		},
	}
}

func TestE(t *testing.T) {

	type test struct {
		args            []interface{}
		expectationFunc func(a *assert.Assertions, actualError *verror)
	}

	cases := map[string]test{
		"op, kind, and subservice in verror": {
			args: []interface{}{
				Op("test-operation"),
				HttpBadRequest,
			},
			expectationFunc: func(a *assert.Assertions, actualError *verror) {
				a.Equal(actualError.Op, Op("test-operation"))
				a.Equal(actualError.Kind, HttpBadRequest)
			},
		},
		"wrapped verror": {
			args: []interface{}{
				Op("wrapping-operation"),
				&verror{
					Op:   Op("test-operation"),
					Kind: HttpBadRequest,
				},
			},
			expectationFunc: func(a *assert.Assertions, actualError *verror) {
				a.Equal(actualError.Op, Op("wrapping-operation"))
				a.Equal(actualError.Kind, HttpBadRequest)
			},
		},
		"create verror from string": {
			args: []interface{}{
				"error-message",
			},
			expectationFunc: func(a *assert.Assertions, actualError *verror) {
				a.Equal(actualError.Err.Error(), "error-message")
			},
		},
		"create verror from error": {
			args: []interface{}{
				errors.New("error-message"),
			},
			expectationFunc: func(a *assert.Assertions, actualError *verror) {
				a.Equal(actualError.Err.Error(), "error-message")
			},
		},
		"add Info, []Info and ErrorInfoProvider to verror": {
			args: []interface{}{
				Info{
					Name: "name1",
					Val:  "val1",
				},
				[]Info{
					{
						Name: "name2",
						Val:  "val2",
					},
				},
				&testErrorInfoProviderPointerReceiver{},
			},
			expectationFunc: func(a *assert.Assertions, actualError *verror) {
				a.Equal(actualError.Infos, []Info{
					{
						Name: "name1",
						Val:  "val1",
					},
					{
						Name: "name2",
						Val:  "val2",
					},
					{
						Name: "outer-info",
						Val: []Info{
							{
								Name: "inner-info",
								Val:  "inner-val",
							},
						},
					},
				})
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			resultError := E(tc.args...)
			tc.expectationFunc(a, resultError.(*verror))

		})
	}
}

func TestIs(t *testing.T) {

	type test struct {
		err      error
		cmprArgs []interface{}
		isEqual  bool
	}

	cases := map[string]test{
		"op, kind, and subservice is equal": {
			err: &verror{
				Op:   "test-operation",
				Kind: HttpBadRequest,
			},
			cmprArgs: []interface{}{
				Op("test-operation"),
				HttpBadRequest,
			},
			isEqual: true,
		},
		"not an *verror": {
			err:      errors.New("cant compare"),
			cmprArgs: []interface{}{},
			isEqual:  false,
		},
		"op is not equal": {
			err: &verror{
				Op:   "test-operation-not-eq",
				Kind: HttpBadRequest,
			},
			cmprArgs: []interface{}{
				Op("test-operation"),
				HttpBadRequest,
			},
			isEqual: false,
		},
		"kind is not equal": {
			err: &verror{
				Op:   "test-operation",
				Kind: HttpClientError,
			},
			cmprArgs: []interface{}{
				Op("test-operation"),
				HttpBadRequest,
			},
			isEqual: false,
		},
		"subservice is not equal": {
			err: &verror{
				Op:   "test-operation",
				Kind: HttpBadRequest,
			},
			cmprArgs: []interface{}{
				Op("test-operation"),
				HttpBadRequest,
			},
			isEqual: false,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			cmprRes := Is(tc.err, tc.cmprArgs...)
			a.Equal(tc.isEqual, cmprRes)
		})
	}
}

func TestError(t *testing.T) {

	type test struct {
		err      *verror
		includes []string
	}

	cases := map[string]test{
		"op, kind, and subservice in Error()": {
			err: &verror{
				Op:   Op("test-operation"),
				Kind: HttpBadRequest,
			},
			includes: []string{
				"test-operation",
				"Okapi HTTP Bad Request Error",
			},
		},
		"info in Error()": {
			err: &verror{
				Infos: []Info{
					{
						Name: "name1",
						Val:  "val1",
					},
					{
						Name: "name2",
						Val:  "val2",
					},
					{
						Name: "outer-info",
						Val: []Info{
							{
								Name: "inner-info",
								Val:  "inner-val",
							},
						},
					},
				},
			},
			includes: []string{
				"name1: val1",
				"name2: val2",
				"outer-info: [{inner-info: inner-val}]",
			},
		},
		"wrapped error in Error()": {
			err: &verror{
				Err: errors.New("actual-error-message"),
			},
			includes: []string{
				"actual-error-message",
			},
		},
		"wrapped verror in Error()": {
			err: &verror{
				Op:  "outer op",
				Err: E(Op("inner op")),
			},
			includes: []string{
				"outer op",
				"inner op",
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			a := assert.New(t)

			errStr := tc.err.Error()

			for _, s := range tc.includes {
				ok := a.Contains(errStr, s)
				if !ok {
					fmt.Println("actual Error string:\n", errStr)
				}
			}
		})
	}
}
