/*
 *  error.go
 *  Created on 29.11.2020
 *  Copyright (C) 2020 Volkswagen AG, All rights reserved.
 */

package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

// The error kind describes the kinds of errors that can exist in VICTOR
type Kind uint16

const (
	Other Kind = iota

	MappingError
	InputError
	SettingsError

	// See usage in repo-extern
	// Is this even anything that can happen?
	DatabaseError
	HttpEmptyOkResponse
	HttpClientError

	HttpNoContent   = Kind(204)
	HttpBadRequest  = Kind(400)
	HttpForbidden   = Kind(403)
	HttpNotFound    = Kind(404)
	HttpInternal    = Kind(500)
	HttpUnavailable = Kind(503)
)

// The type 'SubService' describes the subservice from which the error originates.
type SubService uint8

const (
	_ SubService = iota
	ATB
	REDIS
)

// The ServiceId describes the service in which the error occurred. Its number the first in the error-code created for an error
// set it
type Service uint8

const (
	_ Service = iota
	Sheazuzu
)

var ServiceId Service = 0

// The Error-code is the error code we present to the user of victor in cases where there was an error
func GetErrorCode(err error) int32 {
	verr, ok := err.(*verror)
	if !ok {
		return 0
	}

	return int32(verr.Kind) + int32(verr.SubService)*1000 + int32(ServiceId)*1000000
}

// verror is the fundamental VICTOR-Error-Type
// All internal errors in VICTOR should be of this type
// To create an error which uses this type, use the function E(..)
type verror struct {

	// The Op - shorthand for operation - describes the operation in which an error occurred and was created or has been wrapped
	// It should ideally concisely describe what the service was trying to achieve
	Op Op

	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind

	SubService SubService

	// The underlying error that triggered this one, if any.
	Err error

	// relatively unstructured information about the error
	// examples for things which can be stored in here are features or vehicle-configurations
	// structs which implement the 'ErrorInfoProvider'-interface can be directly passed as an argument to E(...),
	// and will be assigned to the information as described in its implementation of its 'ErrorInfo()'-function
	Infos []Info
}

type Op string

var kindStringMap = map[Kind]string{
	Other:               "Error",
	HttpNoContent:       "HTTP No Content Error",
	HttpBadRequest:      "HTTP Bad Request Error",
	HttpForbidden:       "HTTP Forbidden Error",
	HttpNotFound:        "HTTP Not Found Error",
	HttpInternal:        "HTTP Internal Server Error",
	HttpUnavailable:     "HTTP Service Unavailable Error",
	HttpEmptyOkResponse: "Empty Okapi Response Error",
	HttpClientError:     "HTTP-Client Error",
	MappingError:        "Mapping Error",
	InputError:          "Input Error",
	SettingsError:       "Settings Error",
}

func (k Kind) String() string {
	s, ok := kindStringMap[k]
	if ok {
		return s
	}

	return "Error"
}

var serviceNameMap = map[SubService]string{
	ATB:   "ATB",
	REDIS: "REDIS",
}

func (s SubService) String() string {

	name, ok := serviceNameMap[s]
	if ok {
		return name
	}

	return ""
}

func HttpErrorKindFromStatusCode(status int) Kind {
	return Kind(status)
}

// This map defines which error-kinds result in which http-status-codes when send to the user of vicci
// all undefined Kinds will result in an internal-server-error status-code
var kindHttpStatusMap = map[Kind]int{
	HttpNoContent:  http.StatusNoContent,
	HttpBadRequest: http.StatusBadRequest,
	HttpNotFound:   http.StatusNotFound,
	InputError:     http.StatusBadRequest,
}

// receive the appropriate status-code which to send to the user of VICTOR
func HttpErrorCodeFromError(err error) int {
	verr, ok := err.(*verror)
	if !ok {
		return http.StatusInternalServerError
	}

	s, ok := kindHttpStatusMap[verr.Kind]
	if ok {
		return s
	}

	return http.StatusInternalServerError
}

type Info struct {
	Name string

	Val interface{}
}

func (info Info) String() string {
	return fmt.Sprintf("{%s: %s}", info.Name, info.Val)
}

type ErrorInfoProvider interface {
	ErrorInfo() Info
}

func (err *verror) Error() string {
	sep := "\n\t"

	var b strings.Builder

	fmt.Fprintln(&b)
	printverr(&b, err, sep, true)
	fmt.Fprintln(&b)

	return b.String()
}

func printverr(b *strings.Builder, err error, sep string, isFirstLine bool) {
	if err == nil {
		return
	}

	verr, ok := err.(*verror)

	// In this case we have a usual error, which we regularly want to print and then return
	if !ok {
		fmt.Fprint(b, strings.TrimSpace(err.Error()), sep)
		return
	}

	if isFirstLine {
		fmt.Fprintf(b, "%v %v;\n%s", verr.SubService, verr.Kind, verr.Op)
	} else {
		fmt.Fprint(b, verr.Op)
	}

	if len(verr.Infos) > 0 {
		fmt.Fprint(b, " with:", sep)
	} else {
		fmt.Fprint(b, sep)
	}

	for _, i := range verr.Infos {
		fmt.Fprintf(b, "\t%v: %v%s", i.Name, i.Val, sep)
	}

	printverr(b, verr.Err, sep, false)
}

func E(args ...interface{}) error {

	verr := &verror{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			verr.Op = arg
			// test if op, which is also string falls through

		case SubService:
			verr.SubService = arg

		case *verror:
			verr.Kind = arg.Kind
			verr.SubService = arg.SubService
			verr.Err = arg
			continue
			// test if verr, which is also error falls through

		case error:
			verr.Err = arg

		case Kind:
			verr.Kind = arg

		// TODO(KEM): reduce duplicate infos
		case Info:
			verr.Infos = append(verr.Infos, arg)

		case []Info:
			verr.Infos = append(verr.Infos, arg...)

		case ErrorInfoProvider:
			verr.Infos = append(verr.Infos, arg.ErrorInfo())

		case *ErrorInfoProvider:
			verr.Infos = append(verr.Infos, (*arg).ErrorInfo())

		case string:
			verr.Err = errors.New(arg)
		}
	}

	return verr
}

func Is(actual error, args ...interface{}) bool {

	verr, ok := actual.(*verror)
	if !ok {
		// NOTE(KEM): this function is only for checking Verrors
		return false
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			if verr.Op != arg {
				return false
			}

		case SubService:
			if verr.SubService != arg {
				return false
			}

		case Kind:
			if verr.Kind != arg {
				return false
			}

		case string:
			// TODO(KEM): think about something for this case
			// maybe check if any error contains the string
			// but maybe checking kind and service is enough
		}
	}

	return true
}
