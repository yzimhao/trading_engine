// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: base/variety.proto

package base

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetVarietyRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetVarietyRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetVarietyRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetVarietyRequestMultiError, or nil if none found.
func (m *GetVarietyRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetVarietyRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Symbol

	if len(errors) > 0 {
		return GetVarietyRequestMultiError(errors)
	}

	return nil
}

// GetVarietyRequestMultiError is an error wrapping multiple validation errors
// returned by GetVarietyRequest.ValidateAll() if the designated constraints
// aren't met.
type GetVarietyRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetVarietyRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetVarietyRequestMultiError) AllErrors() []error { return m }

// GetVarietyRequestValidationError is the validation error returned by
// GetVarietyRequest.Validate if the designated constraints aren't met.
type GetVarietyRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetVarietyRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetVarietyRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetVarietyRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetVarietyRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetVarietyRequestValidationError) ErrorName() string {
	return "GetVarietyRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetVarietyRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetVarietyRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetVarietyRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetVarietyRequestValidationError{}

// Validate checks the field values on GetVarietyResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetVarietyResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetVarietyResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetVarietyResponseMultiError, or nil if none found.
func (m *GetVarietyResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetVarietyResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetVarietyResponseMultiError(errors)
	}

	return nil
}

// GetVarietyResponseMultiError is an error wrapping multiple validation errors
// returned by GetVarietyResponse.ValidateAll() if the designated constraints
// aren't met.
type GetVarietyResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetVarietyResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetVarietyResponseMultiError) AllErrors() []error { return m }

// GetVarietyResponseValidationError is the validation error returned by
// GetVarietyResponse.Validate if the designated constraints aren't met.
type GetVarietyResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetVarietyResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetVarietyResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetVarietyResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetVarietyResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetVarietyResponseValidationError) ErrorName() string {
	return "GetVarietyResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetVarietyResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetVarietyResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetVarietyResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetVarietyResponseValidationError{}
