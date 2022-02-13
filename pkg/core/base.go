package core

import (
	"fmt"
	"reflect"
)

func BaseEnv() (*Env, error) {
	env, err := NewEnv(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	// boolean lit
	env.Set("true", Bool(true))
	env.Set("false", Bool(false))
	// nil lit
	env.Set("nil", Nil{})
	// dynamic type
	env.Set("type", Function(func(args []Any) (Any, error) {
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
	env.Set("add", Function(func(args []Any) (Any, error) {
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
	env.Set("mul", Function(func(args []Any) (Any, error) {
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
	env.Set("sub", Function(func(args []Any) (Any, error) {
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
	env.Set("quot", Function(func(args []Any) (Any, error) {
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
	// quotient with default precision
	env.Set("quot*", Function(func(args []Any) (Any, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("wrong number of args for 'quot*': 0")
		}
		result := numberFromInt(1)
		for i, num := range args {
			switch num.(type) {
			default:
				return nil, fmt.Errorf("invalid type for 'quot*': '%v'", num)
			case Number:
				if i == 1 {
					result = args[0].(Number)
				}
				result = result.Quot2(num.(Number))
			}
		}
		return result, nil
	}))
	// remainder with precision
	env.Set("rem", Function(func(args []Any) (Any, error) {
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
		return (args[0].(Number)).Rem(args[1].(Number), args[2].(Number)), nil
	}))
	env.Set("round", Function(func(args []Any) (Any, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("round requires 2 args, got %d", len(args))
		}
		for _, num := range args {
			switch num.(type) {
			default:
				return Zero, fmt.Errorf("invalid type for 'round': '%v'", num)
			case Number:
				break
			}
		}
		return (args[0].(Number)).Round(args[1].(Number)), nil
	}))
	env.Set("prn", Function(func(args []Any) (Any, error) {
		for _, arg := range args {
			fmt.Printf("%v ", arg)
		}
		fmt.Println("")
		return Nil{}, nil
	}))
	env.Set("equal?", Function(func(args []Any) (Any, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("'equal?' expects at least 2 args, got %d", len(args))
		}
		first := args[0]
		for _, item := range args[1:] {
			if !isEqual(first, item) {
				return Bool(false), nil
			}
		}
		return Bool(true), nil
	}))
	return env, nil
}

func isEqual(a Any, b Any) Bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return Bool(false)
	}
	switch a.(type) {
	default:
		return Bool(a == b)
	case List:
		if len(a.(List)) != len(b.(List)) {
			return Bool(false)
		}
		res := true
		for i, itemA := range a.(List) {
			itemB := (b.(List))[i]
			res = res && bool(isEqual(itemA, itemB))
		}
		return Bool(res)
	case Vector:
		if len(a.(Vector)) != len(b.(Vector)) {
			return Bool(false)
		}
		res := true
		for i, itemA := range a.(Vector) {
			itemB := (b.(Vector))[i]
			res = res && bool(isEqual(itemA, itemB))
		}
		return Bool(res)
	case Number:
		return Bool(a.(Number).toDecimal().Equals(b.(Number).toDecimal()))
	}
}
