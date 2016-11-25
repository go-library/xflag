package xflag

import (
	"fmt"
)

type ErrorReason string

const (
	ERR_HELP_REQUEST = "help requested"
	ERR_UNDEFINED    = "undefined flag"
	ERR_EMPTY_VALUE  = "value not provided"
)

type XFlagError struct {
	FlagSet *FlagSet
	Flag    *Flag
	Reason  ErrorReason
	Message string
}

func NewError(fs *FlagSet, flag *Flag, reason ErrorReason, message string) error {
	return &XFlagError{fs, flag, reason, message}
}

func (e *XFlagError) Error() string {
	if e.FlagSet == nil {
		return fmt.Sprintf("XFlagError : %s", e.Reason)
	} else {
		return fmt.Sprintf("XFlagError (%s): %s: %s", e.FlagSet.Name, e.Reason, e.Message)
	}
}

func IsHelpRequest(err error) bool {
	if err, ok := err.(*XFlagError); ok {
		if err.Reason == ERR_HELP_REQUEST {
			return true
		}
	}

	return false
}

func PrintHelp(err error) {
	if err, ok := err.(*XFlagError); ok {
		if err.Reason == ERR_HELP_REQUEST {
			err.FlagSet.PrintHelp()
		}
	}
}
