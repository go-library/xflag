package xflag

import (
	"fmt"
)

type ErrorCode int

const (
	ERROR ErrorCode = iota
	ERROR_HELP_REQUESTED
	ERROR_UNDEFINED_FLAG
	ERROR_EMPTY_VALUE
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

func getError(err error) *Error {
	if err, ok := err.(*Error); ok {
		return err
	}
	return nil
}

func GetErrorCode(err error) ErrorCode {
	if err := getError(err); err != nil {
		return err.Code
	}
	return -1
}

func PrintHelp(err error) {
	if err := getError(err); err != nil {
		err.FlagSet.PrintHelp()
	}
}
