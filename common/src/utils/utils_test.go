/*
 *  utils_test.go
 *  Created on 12.02.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToStringPtr(t *testing.T) {

	type test struct {
		value  string
		result *string
	}

	value := "MyValue"
	cases := map[string]test{
		"simple value": {
			value:  value,
			result: &value,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			assrt := assert.New(t)

			result := ToStringPtr(tc.value)
			assrt.Equal(tc.result, result, "Values do not match")
		})
	}
}

func TestToString(t *testing.T) {

	type test struct {
		value  *string
		result string
	}

	value := "MyValue"
	cases := map[string]test{
		"simple value": {
			value:  &value,
			result: value,
		},
		"nil value": {
			value:  nil,
			result: "",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			assrt := assert.New(t)

			result := ToString(tc.value)
			assrt.Equal(tc.result, result, "Values do not match")
		})
	}
}

func TestToStringPtrOrNil(t *testing.T) {

	type test struct {
		value  string
		result *string
	}

	value := "MyValue"
	cases := map[string]test{
		"simple value": {
			value:  value,
			result: &value,
		},
		"empty value": {
			value:  "",
			result: nil,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			assrt := assert.New(t)

			result := ToStringPtrOrNil(tc.value)
			assrt.Equal(tc.result, result, "Values do not match")
		})
	}
}

func TestToBool(t *testing.T) {

	type test struct {
		value  *bool
		result bool
	}

	value := true
	cases := map[string]test{
		"simple value": {
			value:  &value,
			result: value,
		},
		"nil value": {
			value:  nil,
			result: false,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			assrt := assert.New(t)

			result := ToBool(tc.value)
			assrt.Equal(tc.result, result, "Values do not match")
		})
	}
}

func TestToBoolPtr(t *testing.T) {

	type test struct {
		value  bool
		result *bool
	}

	value := true
	cases := map[string]test{
		"simple value": {
			value:  true,
			result: &value,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			assrt := assert.New(t)

			result := ToBoolPtr(tc.value)
			assrt.Equal(tc.result, result, "Values do not match")
		})
	}
}

func TestToInt64(t *testing.T) {

	type test struct {
		value  *int64
		result int64
	}

	value := int64(1234)
	cases := map[string]test{
		"simple value": {
			value:  &value,
			result: value,
		},
		"nil value": {
			value:  nil,
			result: 0,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {

			assrt := assert.New(t)

			result := ToInt64(tc.value)
			assrt.Equal(tc.result, result, "Values do not match")
		})
	}
}

func TestToStringArray(t *testing.T) {

	t.Parallel()

	type test struct {
		given    *[]string
		expected []string
	}

	cases := map[string]test{
		"regular array": {
			given: &[]string{
				"foo",
				"bar",
			},
			expected: []string{
				"foo",
				"bar",
			},
		},
		"nil": {
			given:    nil,
			expected: nil,
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := ToStringArray(tc.given)
			assert.EqualValues(t, tc.expected, result)
		})
	}
}

func TestToStringArrayPtr(t *testing.T) {

	t.Parallel()

	type test struct {
		given    []string
		expected *[]string
	}

	cases := map[string]test{
		"regular array": {
			given: []string{
				"foo",
				"bar",
			},
			expected: &[]string{
				"foo",
				"bar",
			},
		},
		"nil": {
			given:    nil,
			expected: nil,
		},
		"empty": {
			given:    []string{},
			expected: &[]string{},
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := ToStringArrayPtr(tc.given)
			assert.EqualValues(t, tc.expected, result)
		})
	}
}

func TestToEmptyStringArray(t *testing.T) {

	t.Parallel()

	type test struct {
		given    []string
		expected []string
	}

	cases := map[string]test{
		"non empty list": {
			given:    []string{"foo", "bar"},
			expected: []string{"foo", "bar"},
		},
		"empty list": {
			given:    []string{},
			expected: []string{},
		},
		"nil list": {
			given:    nil,
			expected: []string{},
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := ToEmptyStringArray(tc.given)
			assert.EqualValues(t, tc.expected, result)
		})
	}
}

func TestDeleteEmpty(t *testing.T) {

	t.Parallel()

	type test struct {
		input    []string
		expected []string
	}

	cases := map[string]test{
		"nil array": {
			input:    nil,
			expected: nil,
		},
		"empty array": {
			input:    []string{},
			expected: nil,
		},
		"array with one empty element": {
			input:    []string{""},
			expected: nil,
		},
		"array with elements and one empty element": {
			input:    []string{"a", "", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for name, tc := range cases {

		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := DeleteEmpty(tc.input)
			assert.EqualValues(t, tc.expected, result)
		})
	}
}
