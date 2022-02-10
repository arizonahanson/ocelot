package core

import (
	"strconv"
)

type String string

func (value String) ToNumber() (*Number, error) {
	return numberFromString(value)
}

func (value String) Unquote() (*String, error) {
	unquoted, err := strconv.Unquote(string(value))
	if err != nil {
		return nil, err
	}
	s := String(unquoted)
	return &s, nil
}
