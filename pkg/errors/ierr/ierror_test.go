package ierr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Code(t *testing.T) {
	type fields struct {
		code Code
	}
	tests := []struct {
		name   string
		fields fields
		want   Code
	}{
		{
			name:   "code canceled",
			want:   Canceled,
			fields: fields{code: Canceled},
		},
		{
			name:   "code internal",
			want:   Internal,
			fields: fields{code: Internal},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{code: tt.fields.code}
			got := e.Code()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestError_Message(t *testing.T) {
	type fields struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "error message",
			want:   "error message",
			fields: fields{msg: "error message"},
		},
		{
			name:   "error message next",
			want:   "error message next",
			fields: fields{msg: "error message next"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{msg: tt.fields.msg}
			assert.Equal(t, tt.want, e.Message())
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	type fields struct {
		original error
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name:   "original error",
			fields: fields{original: errors.New("original error")},
			want:   errors.New("original error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{original: tt.fields.original}
			assert.Equal(t, tt.want, e.Unwrap())
		})
	}
}

func TestGetCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "not found code",
			args: args{err: errors.New("not found code")},
			want: Unknown,
		},
		{
			name: "internal error",
			args: args{err: New(Internal, "internal error")},
			want: Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetCode(tt.args.err), "GetCode(%v)", tt.args.err)
		})
	}
}

func TestGetMessage(t *testing.T) {
	var errNil *Error
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty message, error not a pointer",
			args: args{err: Error{}},
			want: "",
		},
		{
			name: "typed nil",
			args: args{err: errNil},
			want: "",
		},
		{
			name: "message in standard error",
			args: args{err: errors.New("message in standard error")},
			want: "message in standard error",
		},
		{
			name: "message in internal error",
			args: args{err: New(Internal, "message in internal error")},
			want: "message in internal error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetMessage(tt.args.err), "GetMessage(%v)", tt.args.err)
		})
	}
}

func TestError_Error(t *testing.T) {
	st := stack{}
	type fields struct {
		msg      string
		original error
		code     Code
		stack    *stack
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "without original error",
			fields: fields{
				msg:   "Internal error",
				code:  Internal,
				stack: &st,
			},
			want: "code: INTERNAL error: Internal error",
		},
		{
			name: "with original error",
			fields: fields{
				msg:      "Internal error",
				code:     Internal,
				stack:    &st,
				original: errors.New("original error"),
			},
			want: "code: INTERNAL error: Internal error: original error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				msg:      tt.fields.msg,
				original: tt.fields.original,
				code:     tt.fields.code,
				stack:    tt.fields.stack,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}

func TestNewUnknown(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "unknown error",
			args: args{message: "unknown error"},
			want: &Error{
				msg:  "unknown error",
				code: Unknown,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewUnknown(tt.args.message)
			assert.Equal(t, tt.want.msg, e.msg)
			assert.Equal(t, tt.want.code, e.code)
		})
	}
}

func TestErrorf(t *testing.T) {
	type args struct {
		code   Code
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "formatted error",
			args: args{
				format: "%s error",
				code:   Unknown,
				args:   []interface{}{"formatted"},
			},
			want: &Error{
				msg:  "formatted error",
				code: Unknown,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Errorf(tt.args.code, tt.args.format, tt.args.args...)
			assert.Equal(t, tt.want.msg, e.msg)
			assert.Equal(t, tt.want.code, e.code)
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		code     Code
		original error
		message  string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "wrapped error",
			args: args{
				code:     Unknown,
				original: errors.New("original error"),
				message:  "wrapped error",
			},
			want: &Error{
				msg:      "wrapped error",
				original: errors.New("original error"),
				code:     Unknown,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := WrapCode(tt.args.code, tt.args.original, tt.args.message)
			assert.Equal(t, tt.want.msg, e.msg)
			assert.Equal(t, tt.want.code, e.code)
			assert.Equal(t, tt.want.original, e.original)
		})
	}
}

func TestWrapf(t *testing.T) {
	type args struct {
		code     Code
		original error
		format   string
		args     []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "wrapped error",
			args: args{
				code:     Unknown,
				original: errors.New("original error"),
				format:   "%v error",
				args:     []interface{}{"unknown"},
			},
			want: &Error{
				msg:      "unknown error",
				original: errors.New("original error"),
				code:     Unknown,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := WrapfCode(tt.args.code, tt.args.original, tt.args.format, tt.args.args...)
			assert.Equal(t, tt.want.msg, e.msg)
			assert.Equal(t, tt.want.code, e.code)
			assert.Equal(t, tt.want.original, e.original)
		})
	}
}

func TestCode_String(t *testing.T) {
	tests := []struct {
		name string
		c    Code
		want string
	}{
		{name: "Unknown", c: Unknown, want: "UNKNOWN"},
		{name: "Canceled", c: Canceled, want: "CANCELED"},
		{name: "InvalidArgument", c: InvalidArgument, want: "INVALID_ARGUMENT"},
		{name: "NotFound", c: NotFound, want: "NOT_FOUND"},
		{name: "AlreadyExists", c: AlreadyExists, want: "ALREADY_EXISTS"},
		{name: "PermissionDenied", c: PermissionDenied, want: "PERMISSION_DENIED"},
		{name: "ResourceExhausted", c: ResourceExhausted, want: "RESOURCE_EXHAUSTED"},
		{name: "Aborted", c: Aborted, want: "ABORTED"},
		{name: "OutOfRange", c: OutOfRange, want: "OUT_OF_RANGE"},
		{name: "Unimplemented", c: Unimplemented, want: "UNIMPLEMENTED"},
		{name: "Internal", c: Internal, want: "INTERNAL"},
		{name: "Unavailable", c: Unavailable, want: "UNAVAILABLE"},
		{name: "DataLoss", c: DataLoss, want: "DATA_LOSS"},
		{name: "Unauthenticated", c: Unauthenticated, want: "UNAUTHENTICATED"},
		{name: "UnprocessableEntity", c: UnprocessableEntity, want: "UNPROCESSABLE_ENTITY"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.String())
		})
	}
}
