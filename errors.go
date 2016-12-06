package xflag

import (
	"fmt"
)

type ErrorCode uint

const (
	HELP_REQUESTED ErrorCode = iota
	PARSE_ERROR_UNDEFINED_FLAG
	PARSE_ERROR_EMPTY_VALUE
)

type Error struct {
	error
	Code    ErrorCode
	FlagSet *FlagSet
	Flag    *Flag
}

func Errorf(fs *FlagSet, flag *Flag, code ErrorCode, format string, args ...interface{}) *Error {
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

func GetError(err error) *Error {
	if err, ok := err.(*Error); ok {
		return err
	}
	return nil
}

func GetErrorCode(err error) ErrorCode {
	if err := GetError(err); err != nil {
		return err.Code
	}
	return 0
}

func PrintHelp(err error) {
	if err := GetError(err); err != nil {
		err.FlagSet.PrintHelp()
	}
}

func IsHelpRequest(err error) bool {
	if GetErrorCode(err) == HELP_REQUESTED {
		return true
	}
	return false
}
