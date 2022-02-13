package core

import (
	"fmt"
	"strconv"
)

type Any interface{}

type Nil struct{}

func (value Nil) String() string {
	return "nil"
}

type List []Any

func (value List) String() string {
	result := "("
	for i, item := range value {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%v", item)
	}
	return result + ")"
}

type Symbol string

type Bool bool

func (value Bool) String() string {
	if value {
		return "true"
	} else {
		return "false"
	}
}

type String string

func (value String) Number() (Number, error) {
	return numberFromString(value)
}

func (value String) Unquote() (String, error) {
	unquoted, err := strconv.Unquote(string(value))
	if err != nil {
		return "", err
	}
	s := String(unquoted)
	return s, nil
}

type Vector []Any

type Function func(args ...Any) (Any, error)
