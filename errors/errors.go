/*
 *  Dogg3rz is a decentralized metadata version control system
 *  Copyright (C) 2019 D. Andrew Padilla dba Datacequia
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as
 *  published by the Free Software Foundation, either version 3 of the
 *  License, or (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU Affero General Public License for more details.
 *
 *  You should have received a copy of the GNU Affero General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package errors

import (
	"fmt"

	"encoding/json"

	"github.com/pkg/errors"
)

type ErrorType uint

const (
	NoType = ErrorType(iota)
	InvalidArg
	OutOfRange
	InvalidPathElement
	AlreadyExists
	NotFound
	UnexpectedType
	InitNodeError
	TryAgain
	ConfigError
	RepoError
	UnexpectedValue
	AssertionError
)

type badDogg3rz struct {
	errorType     ErrorType
	originalError error
	contextInfo   map[string]string
	//context errorContext
}

/*
type errorContext struct {
	Field   string
	Message string
}
*/

func errorTypeToString(errType ErrorType) string {

	switch errType {
	case NoType:
		return "external error"
	case InvalidArg:
		return "invalid argument"
	case OutOfRange:
		return "out of range"
	case InvalidPathElement:
		return "invalid dogg3rz path element value"
	case AlreadyExists:
		return "resource already exists"
	case NotFound:
		return "resource not found"
	case UnexpectedType:
		return "unexpected type encountered"
	case InitNodeError:
		return "node initialization error"
	case TryAgain:
		return "try again"
	case ConfigError:
		return "configuration error"
	case RepoError:
		return "repository error"
	case UnexpectedValue:
		return "unexpected value encountered"
	case AssertionError:
		return "assertion error"
	default:
		return "unknown error"
	}
}

// Error returns the mssage of a customError
func (error badDogg3rz) Error() string {

	if len(error.contextInfo) > 0 {
		context, _ := json.Marshal(error.contextInfo)
		return fmt.Sprintf("%s - %s:  %v", errorTypeToString(error.errorType), string(context), error.originalError)
	}

	return fmt.Sprintf("%s: %v", errorTypeToString(error.errorType), error.originalError)
}

// New creates a new customError
func (errType ErrorType) New(msg string) error {
	return badDogg3rz{errorType: errType, originalError: errors.New(msg), contextInfo: make(map[string]string)}
}

// New creates a new customError with formatted message
func (errType ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)

	return badDogg3rz{errorType: errType, originalError: err, contextInfo: make(map[string]string)}
}

// Wrap creates a new wrapped error
func (errType ErrorType) Wrap(err error, msg string) error {
	return errType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return badDogg3rz{errorType: errType, originalError: newErr}
}

// New creates a no type error
func New(msg string) error {
	return badDogg3rz{errorType: NoType, originalError: errors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return badDogg3rz{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap wrans an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(badDogg3rz); ok {
		return badDogg3rz{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contextInfo:   customErr.contextInfo,
		}
	}

	return badDogg3rz{errorType: NoType, originalError: wrappedError}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field string, message string) error {
	//	context := errorContext{Field: field, Message: message}
	if myErr, ok := err.(badDogg3rz); ok {
		myErr.contextInfo[field] = message

		return myErr //return badDogg3rz{errorType: myErr.errorType, originalError: myErr.originalError, context: context}
	}

	return badDogg3rz{errorType: NoType, originalError: err, contextInfo: make(map[string]string)}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	//	emptyContext := errorContext{}
	if myErr, ok := err.(badDogg3rz); ok {

		//return map[string]string{"field": myErr.context.Field, "message": myErr.context.Message}
		return myErr.contextInfo
	}

	return make(map[string]string)
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if myErr, ok := err.(badDogg3rz); ok {
		return myErr.errorType
	}

	return NoType
}
