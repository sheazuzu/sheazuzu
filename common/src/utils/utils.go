/*
 *  utils.go
 *  Created on 12.02.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package utils

import (
	"encoding/json"
)

func ToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func ToStringArray(value *[]string) []string {
	if value == nil {
		return nil
	}
	return *value
}

func ToStringArrayPtr(arr []string) *[]string {

	if arr == nil {
		return nil
	}

	return &arr
}

func ToStringPtr(value string) *string {
	return &value
}

func ToStringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func ToBool(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func ToBoolPtr(value bool) *bool {
	return &value
}

func ToInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func ToInt32(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}

func ToInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func ToIntPtr(value int) *int {
	return &value
}

func ToInt32Ptr(value int32) *int32 {
	return &value
}

func ToFloat32(value *float32) float32 {
	if value == nil {
		return 0.0
	}
	return *value
}

func ToFloat32Ptr(value float32) *float32 {
	return &value
}

func ToFloat64Ptr(value float64) *float64 {
	return &value
}

func ToEmptyStringArray(arr []string) []string {

	if arr == nil {
		return make([]string, 0)
	}

	return arr
}

// Copy the in value by marshaling/unmarshalling the source value
func DeepCopy(in, out interface{}) {

	bytes, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		panic(err)
	}
}

func DeleteEmpty(arr []string) []string {
	var result []string
	for _, str := range arr {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func MinIntArray(l []int) int {
	min := l[0]
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return min
}
