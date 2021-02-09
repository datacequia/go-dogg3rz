/*
 * Copyright (c) 2019-2020 Datacequia LLC. All rights reserved.
 *
 * This program is licensed to you under the Apache License Version 2.0,
 * and you may not use this file except in compliance with the Apache License Version 2.0.
 * You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the Apache License Version 2.0 is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
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
	InvalidValue
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
	UnhandledValue
	InvalidState      // ATTEMPT TO EXECUTE OPERATION IN AN INVALID STATE
	Cancelled         // OPERATION WAS CANCELLED
	TimedOut          // TIMEOUT OCCURRED WHILE WAITING TO PERFORM SOME OPERATION
	RollbackRequested // USER ISSUED A ROLLBACK
	ChannelClosed     // A CLOSED CHANNEL WAS DETECTED
	EmptyCommit       // AN ATTEMPT TO COMMIT A RESOURCE BUT NOTHING TO COMMIT
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
		return "NoType"
	case InvalidValue:
		return "InvalidValue"
	case OutOfRange:
		return "OutOfRange"
	case InvalidPathElement:
		return "InvalidPathElement"
	case AlreadyExists:
		return "AlreadyExists"
	case NotFound:
		return "NotFound"
	case UnexpectedType:
		return "UnexpectedType"
	case InitNodeError:
		return "InitNodeError"
	case TryAgain:
		return "TryAgain"
	case ConfigError:
		return "ConfigError"
	case RepoError:
		return "RepoError"
	case UnexpectedValue:
		return "UnexpectedValue"
	case AssertionError:
		return "AssertionError"
	case UnhandledValue:
		return "UnhandledValue"
	case InvalidState:
		return "InvalidState"
	case Cancelled:
		return "Cancelled"
	case TimedOut:
		return "TimedOut"
	case RollbackRequested:
		return "RollbackRequested"
	case ChannelClosed:
		return "ChannelClosed"
	case EmptyCommit:
		return "EmptyCommit"
	default:
		panic("unknown error")
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
