package gexpect

import (
	"errors"
)

type ValueNotBindError struct {
	VarName string
}

func (e ValueNotBindError) Error() string {
	return "Value not bind: " + e.VarName
}

type ValueNotFoundError struct {
	VarName string
}

func (e ValueNotFoundError) Error() string {
	return "Value not found in list: " + e.VarName
}

type TerminatedError struct {
	Message string
}

func (e TerminatedError) Error() string {
	return "flow terminated with message: " + e.Message
}

var (
	EOF     = errors.New("EOF")
	TIMEOUT = errors.New("Timeout")
)
