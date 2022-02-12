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
			return nil, fmt.Errorf("too many args for 'type': %d", len(args))
		}
		typeStr := reflect.TypeOf(args[0]).String()
		return String(typeStr), nil
	})
	e[Symbol("add")] = Function(func(args ...Any) (Any, error) {
		result := Zero
		for _, num := range args {
			switch num.(type) {
			default:
				return result, fmt.Errorf("invalid type for 'add': '%v'", num)
			case Number:
				result = result.Add(num.(Number))
				break
			}
		}
		return result, nil
	})
	e[Symbol("mul")] = Function(func(args ...Any) (Any, error) {
		result := numberFromInt(1)
		for _, num := range args {
			switch num.(type) {
			default:
				return Zero, fmt.Errorf("invalid type for 'mul': '%v'", num)
			case Number:
				result = result.Mul(num.(Number))
				break
			}
		}
		return result, nil
	})
	e[Symbol("sub")] = Function(func(args ...Any) (Any, error) {
		result := Zero
		for i, num := range args {
			switch num.(type) {
			default:
				return Zero, fmt.Errorf("invalid type for 'sub': '%v'", num)
			case Number:
				if i == 1 {
					result = args[0].(Number)
				}
				result = result.Sub(num.(Number))
				break
			}
		}
		return result, nil
	})
	e[Symbol("quot")] = Function(func(args ...Any) (Any, error) {
		if len(args) != 3 {
			return Zero, fmt.Errorf("quot requires 3 args, got %d", len(args))
		}
		for _, num := range args {
			switch num.(type) {
			default:
				return Zero, fmt.Errorf("invalid type for 'quot': '%v'", num)
			case Number:
				break
			}
		}
		return (args[0].(Number)).Quot((args[1].(Number)), (args[2].(Number))), nil
	})
	e[Symbol("rem")] = Function(func(args ...Any) (Any, error) {
		if len(args) != 3 {
			return Zero, fmt.Errorf("rem requires 3 args, got %d", len(args))
		}
		for _, num := range args {
			switch num.(type) {
			default:
				return Zero, fmt.Errorf("invalid type for 'rem': '%v'", num)
			case Number:
				break
			}
		}
		return (args[0].(Number)).Rem((args[1].(Number)), (args[2].(Number))), nil
	})
	return e
}
