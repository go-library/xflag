package xflag

import (
	"fmt"
)

type ErrorCode uint

const (
	ERR_HELP_REQUESTED ErrorCode = iota
	ERR_UNDEFINED
	ERR_EMPTY_VALUE
)

type Error struct {
	error
	Code    ErrorCode
	FlagSet *FlagSet
	Flag    *Flag
}

func Errorf(fs *FlagSet, flag *Flag, code ErrorCode, format string, args ...interface{}) error {
	return &Error{
		error:   fmt.Errorf(format, args...),
		Code:    code,
		FlagSet: fs,
		Flag:    flag,
	}
}

func (e *Error) Error() string {
	switch {
	case e.FlagSet != nil && e.Flag != nil:
		return fmt.Sprintf("%s:%s: %s", e.FlagSet, e.Flag, e.error.Error())
	case e.FlagSet != nil:
		return fmt.Sprintf("%s: %s", e.FlagSet, e.error.Error())
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

func PrintHelp(err error) {
	if err, ok := err.(*Error); ok {
		if err.Code == ERR_HELP_REQUESTED {
			err.FlagSet.PrintHelp()
		}
	}
}

func IsFlagError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

func GetFlagError(err error) *Error {
	if err, ok := err.(*Error); ok {
		return err
	}
	return nil
}
