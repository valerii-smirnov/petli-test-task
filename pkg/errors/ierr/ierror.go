package ierr

import (
	"fmt"
	"reflect"
)

// KV is an alias for building generic key/value
type KV map[string]string

// Code is an unsigned 32-bit error code as defined for internal errors.
type Code uint32

// Status codes for internal errors.
const (
	Unknown             Code = 0
	Canceled            Code = 1
	InvalidArgument     Code = 2
	NotFound            Code = 3
	AlreadyExists       Code = 4
	PermissionDenied    Code = 5
	ResourceExhausted   Code = 6
	Aborted             Code = 7
	OutOfRange          Code = 8
	Unimplemented       Code = 9
	Internal            Code = 10
	Unavailable         Code = 11
	DataLoss            Code = 12
	Unauthenticated     Code = 13
	UnprocessableEntity Code = 14
	DeadlineExceeded    Code = 15
	FailedPrecondition  Code = 16
)

const (
	unknownKey                  = "UNKNOWN"
	canceledKey                 = "CANCELED"
	invalidArgumentKey          = "INVALID_ARGUMENT"
	notFoundKey                 = "NOT_FOUND"
	alreadyExistsKey            = "ALREADY_EXISTS"
	permissionDeniedKey         = "PERMISSION_DENIED"
	resourceExhaustedKey        = "RESOURCE_EXHAUSTED"
	abortedKey                  = "ABORTED"
	outOfRangeKey               = "OUT_OF_RANGE"
	unimplementedKey            = "UNIMPLEMENTED"
	internalKey                 = "INTERNAL"
	unavailableKey              = "UNAVAILABLE"
	dataLossKey                 = "DATA_LOSS"
	unauthenticatedKey          = "UNAUTHENTICATED"
	unprocessableEntityKey      = "UNPROCESSABLE_ENTITY"
	deadlineExceededEntityKey   = "DEADLINE_EXCEEDED"
	failedPreconditionEntityKey = "FAILED_PRECONDITION"
)

// String - represents code in string code key.
func (c Code) String() string {
	switch c {
	case Unknown:
		return unknownKey
	case Canceled:
		return canceledKey
	case InvalidArgument:
		return invalidArgumentKey
	case NotFound:
		return notFoundKey
	case AlreadyExists:
		return alreadyExistsKey
	case PermissionDenied:
		return permissionDeniedKey
	case ResourceExhausted:
		return resourceExhaustedKey
	case Aborted:
		return abortedKey
	case OutOfRange:
		return outOfRangeKey
	case Unimplemented:
		return unimplementedKey
	case Internal:
		return internalKey
	case Unavailable:
		return unavailableKey
	case DataLoss:
		return dataLossKey
	case Unauthenticated:
		return unauthenticatedKey
	case UnprocessableEntity:
		return unprocessableEntityKey
	case DeadlineExceeded:
		return deadlineExceededEntityKey
	case FailedPrecondition:
		return failedPreconditionEntityKey
	default:
		return unknownKey
	}
}

// New - creates new instance of internal error.
func New(code Code, message string) *Error {
	return &Error{
		msg:   message,
		code:  code,
		stack: callers(),
	}
}

// NewUnknown - creates new instance of internal error with UNKNOWN status code.
func NewUnknown(message string) *Error {
	return &Error{
		msg:   message,
		code:  Unknown,
		stack: callers(),
	}
}

// Errorf - creates new instance of internal error including Code with string formatting.
func Errorf(code Code, format string, args ...interface{}) *Error {
	return &Error{
		msg:   fmt.Sprintf(format, args...),
		code:  code,
		stack: callers(),
	}
}

// WrapCode - creates new instance of internal error including Code with error original.
func WrapCode(code Code, original error, message string) *Error {
	return &Error{
		msg:      message,
		code:     code,
		original: original,
		stack:    callers(),
	}
}

// Wrap - creates new instance of internal error with error original.
// if original error doesn't contain Code (ierr), Unknown will be set
func Wrap(original error, message string) *Error {
	return &Error{
		msg:      message,
		code:     GetCode(original),
		original: original,
		stack:    callers(),
	}
}

// WrapfCode - creates new instance of internal error with error including Code original and string formation.
func WrapfCode(code Code, original error, format string, args ...interface{}) *Error {
	return &Error{
		msg:      fmt.Sprintf(format, args...),
		code:     code,
		original: original,
		stack:    callers(),
	}
}

// Wrapf - creates new instance of internal error with error original and string formation.
// if original error doesn't contain Code (ierr), Unknown will be set
func Wrapf(original error, format string, args ...interface{}) *Error {
	return &Error{
		msg:      fmt.Sprintf(format, args...),
		code:     GetCode(original),
		original: original,
		stack:    callers(),
	}
}

// Error - represents internal error.
type Error struct {
	msg        string
	original   error
	code       Code
	properties KV
	stack      *stack
}

// Error - string representation of error.
func (e Error) Error() string {
	stack := e.stack.String()
	code := e.code.String()

	if e.original == nil {
		return fmt.Sprintf("code: %s error: %s%s", code, e.msg, stack)
	}

	originMessage := e.original.Error()
	if se, ok := e.original.(*Error); ok {
		originMessage = se.Message()
		stack = se.stack.String()
	}

	msg := e.msg + ": " + originMessage

	return fmt.Sprintf("code: %s error: %s%s", code, msg, stack)
}

// Code - returns error code.
func (e Error) Code() Code { return e.code }

// Unwrap - returns error original.
func (e Error) Unwrap() error { return e.original }

// Message - returns error msg.
func (e Error) Message() string { return e.msg }

// OriginalError - returns original error
func (e Error) OriginalError() error { return e.original }

// Props - returns properties error
func (e Error) Props() KV { return e.properties }

// SetProps - set additional properties if necessary
func (e *Error) SetProps(kv KV) *Error {
	e.properties = kv

	return e
}

// GetCode - returns code from error.
func GetCode(err error) Code {
	if se, ok := err.(*Error); ok {
		return se.Code()
	}

	return Unknown
}

// GetMessage - return internal message without stack.
func GetMessage(err error) string {
	if reflect.ValueOf(err).Kind() != reflect.Ptr {
		return ""
	}

	if reflect.ValueOf(err).Kind() == reflect.Ptr && reflect.ValueOf(err).IsNil() {
		return ""
	}

	if se, ok := err.(*Error); ok {
		return se.Message()
	}

	return err.Error()
}

// GetProps - return properties of internal error if exists
func GetProps(err error) KV {
	if se, ok := err.(*Error); ok {
		return se.Props()
	}

	return nil
}
