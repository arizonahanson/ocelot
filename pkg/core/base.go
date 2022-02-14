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
		result := numberFromInt(0)
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
				return numberFromInt(0), fmt.Errorf("invalid type for 'mul': '%v'", num)
			case Number:
				result = result.Mul(num.(Number))
				break
			}
		}
		return result, nil
	}))
	// subtraction
	env.Set("sub", Function(func(args []Any) (Any, error) {
		result := numberFromInt(0)
		for i, num := range args {
			switch num.(type) {
			default:
				return numberFromInt(0), fmt.Errorf("invalid type for 'sub': '%v'", num)
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
			return numberFromInt(0), fmt.Errorf("quot requires 3 args, got %d", len(args))
		}
		for _, num := range args {
			switch num.(type) {
			default:
				return numberFromInt(0), fmt.Errorf("invalid type for 'quot': '%v'", num)
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
			return numberFromInt(0), fmt.Errorf("rem requires 3 args, got %d", len(args))
		}
		for _, num := range args {
			switch num.(type) {
			default:
				return numberFromInt(0), fmt.Errorf("invalid type for 'rem': '%v'", num)
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
				return numberFromInt(0), fmt.Errorf("invalid type for 'round': '%v'", num)
			case Number:
				break
			}
		}
		return (args[0].(Number)).Round(args[1].(Number)), nil
	}))
	env.Set("prn", Function(func(args []Any) (Any, error) {
		fmt.Println(List(args))
		return Nil{}, nil
	}))
	env.Set("eq?", Function(func(args []Any) (Any, error) {
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
	env.Set("not", Function(func(args []Any) (Any, error) {
		if len(args) > 1 {
			return nil, fmt.Errorf("'not' expects one arg, got %d", len(args))
		}
		return Bool(!IsTruthy(args[0])), nil
	}))
	env.Set("lt?", Function(func(args []Any) (Any, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("'lt?' expects two args, got %d", len(args))
		}
		if !isNumber(args[0]) || !isNumber(args[1]) {
			return nil, fmt.Errorf("'lt?' only works with numbers: %s, %s", args[0], args[1])
		}
		return Bool(args[0].(Number).toDecimal().LessThan(args[1].(Number).toDecimal())), nil
	}))
	env.Set("lteq?", Function(func(args []Any) (Any, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("'lteq?' expects two args, got %d", len(args))
		}
		if !isNumber(args[0]) || !isNumber(args[1]) {
			return nil, fmt.Errorf("'lteq?' only works with numbers: %s, %s", args[0], args[1])
		}
		return Bool(args[0].(Number).toDecimal().LessThanOrEqual(args[1].(Number).toDecimal())), nil
	}))
	env.Set("gt?", Function(func(args []Any) (Any, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("'gt?' expects two args, got %d", len(args))
		}
		if !isNumber(args[0]) || !isNumber(args[1]) {
			return nil, fmt.Errorf("'gt?' only works with numbers: %s, %s", args[0], args[1])
		}
		return Bool(args[0].(Number).toDecimal().GreaterThan(args[1].(Number).toDecimal())), nil
	}))
	env.Set("gteq?", Function(func(args []Any) (Any, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("'gteq?' expects two args, got %d", len(args))
		}
		if !isNumber(args[0]) || !isNumber(args[1]) {
			return nil, fmt.Errorf("'gteq?' only works with numbers: %s, %s", args[0], args[1])
		}
		return Bool(args[0].(Number).toDecimal().GreaterThanOrEqual(args[1].(Number).toDecimal())), nil
	}))
	return env, nil
}

func isNumber(any Any) bool {
	switch any.(type) {
	default:
		return false
	case Number:
		return true
	}
}

func IsTruthy(any Any) bool {
	switch any.(type) {
	default:
		return true
	case Bool:
		return bool(any.(Bool))
	case Nil:
		return false
	}
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
