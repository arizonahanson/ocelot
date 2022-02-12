package core

import (
	"fmt"
	"reflect"
)

func GetEnv(env map[string]interface{}) Environment {
	e := make(Environment)
	e[Symbol("true")] = Bool(true)
	e[Symbol("false")] = Bool(false)
	e[Symbol("nil")] = Nil{}
	e[Symbol("type")] = Function(func(args ...Any) (Any, error) {
		if len(args) == 0 {
			return Nil{}, nil
		}
		if len(args) > 1 {
			return nil, fmt.Errorf("too many args for type: %d", len(args))
		}
		typeStr := reflect.TypeOf(args[0]).String()
		return String(typeStr), nil
	})
	e[Symbol("add")] = Function(func(args ...Any) (Any, error) {
		result := Zero
		for _, num := range args {
			switch num.(type) {
			default:
				return result, fmt.Errorf("invalid type for add: %v", num)
			case Number:
				result = result.Add(num.(Number))
				break
			}
		}
		return result, nil
	})
	return e
}
