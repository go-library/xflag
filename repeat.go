package xflag

import (
	"fmt"
)

type StringList []string

func NewStringList() *StringList {
	return &StringList{}
}

func (l *StringList) Set(s string) error {
	*l = append(*l, s)
	return nil
}

func (l *StringList) String() string {
	return fmt.Sprintf("%#v", ([]string)(*l))
}
