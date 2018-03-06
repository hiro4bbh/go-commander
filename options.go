package gocommander

import (
	"fmt"
	"strings"
)

// Option is the interface for option values.
type Option interface {
	// Set sets the value parsed from the given str, and returns error if occurred.
	// str has the form "", "+", "-", or "=$VALUE".
	Set(str string) error
	// String returns the string representation of the option value.
	String() string
	// ValueFormat returns the option value format.
	ValueFormat() string
}

// OptionBool is Option having a bool variable.
type OptionBool struct {
	value bool
}

// NewOptionBool returns a new OptionBool with the given default value.
func NewOptionBool(value bool) *OptionBool {
	return &OptionBool{value}
}

// Get returns the value.
func (opt *OptionBool) Get() bool {
	return opt.value
}

// Set is for interface Option.
func (opt *OptionBool) Set(str string) error {
	switch strings.ToLower(str) {
	case "-":
		opt.value = false
	case "", "+":
		opt.value = true
	default:
		return fmt.Errorf("illegal OptionBool value: %s", str)
	}
	return nil
}

// String is for interface Option.
func (opt *OptionBool) String() string {
	return fmt.Sprintf("%v", opt.value)
}

// ValueFormat is for interface Option.
func (opt *OptionBool) ValueFormat() string {
	return "[+-]"
}

// OptionString is Option having a string variable.
type OptionString struct {
	value string
}

// NewOptionString returns a new OptionString with the given default value.
func NewOptionString(value string) *OptionString {
	return &OptionString{value}
}

// Get returns the value.
func (opt *OptionString) Get() string {
	return opt.value
}

// Set is for interface Option.
func (opt *OptionString) Set(str string) error {
	if str == "" {
		opt.value = ""
	} else if strings.HasPrefix(str, "=") {
		opt.value = str[len("="):]
	} else {
		return fmt.Errorf("illegal OptionString value: %s", str)
	}
	return nil
}

// String is for interface Option.
func (opt *OptionString) String() string {
	return fmt.Sprintf("%q", opt.value)
}

// ValueFormat is for interface Option.
func (opt *OptionString) ValueFormat() string {
	return "=VALUE"
}
