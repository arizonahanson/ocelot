package core

import (
	"fmt"
	"reflect"
)

func BaseEnv() Env {
	env := NewEnv(nil)
	// boolean lit
	env.Set("true", Bool(true))
	env.Set("false", Bool(true))
	// nil lit
	env.Set("nil", Nil{})
	// dynamic type
	env.Set("type", Function(func(args ...Any) (Any, error) {
		if len(args) == 0 {
			return Nil{}, nil
		}
		if len(args) > 1 {
			return nil, fmt.Errorf("too many args for 'type': %d", len(args))
		}
		typeStr := reflect.TypeOf(args[0]).String()
		return String(typeStr), nil
	}))
	// addition
	env.Set("add", Function(func(args ...Any) (Any, error) {
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
	}))
	// multiplication
	env.Set("mul", Function(func(args ...Any) (Any, error) {
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
	}))
	// subtraction
	env.Set("sub", Function(func(args ...Any) (Any, error) {
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
	}))
	// quotient with precision
	env.Set("quot", Function(func(args ...Any) (Any, error) {
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
	}))
	// remainder with precision
	env.Set("rem", Function(func(args ...Any) (Any, error) {
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
	}))
	return env
}
