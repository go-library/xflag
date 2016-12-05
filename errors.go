package xflag

import (
	"fmt"
)

type ErrorCode uint

type Error struct {
	error
	Code    ErrorCode
	FlagSet *FlagSet
	Flag    *Flag
}

type ErrorReason string

const (
	ERR_HELP_REQUESTED ErrorCode = iota
	ERR_UNDEFINED
	ERR_EMPTY_VALUE
)

func NewError(fs *FlagSet, flag *Flag, code ErrorCode, err error) error {
	return &Error{
		error:   err,
		Code:    code,
		FlagSet: fs,
		Flag:    flag,
	}
}

func (e *Error) Error() string {
	switch {
	case e.FlagSet != nil && e.Flag != nil:
		return fmt.Sprintf("(%s(%s): %s", e.FlagSet.Name, e.Flag, e.Error())
	case e.FlagSet != nil:
		return fmt.Sprintf("(%s: %s", e.FlagSet.Name, e.Error())
	default:
		return e.Error()
	}
}

func IsHelpRequest(err error) bool {
	if err, ok := err.(*Error); ok {
		if err.Code == ERR_HELP_REQUESTED {
			return true
		}
	}
	return false
}
