package base

import "github.com/starlight/ocelot/pkg/core"

type ThunkType func() (core.Any, error)

type EvalType func(ast core.Any, env core.Env) (core.Any, error)

func (eval EvalType) Thunk(ast core.Any, env core.Env) ThunkType {
	return func() (core.Any, error) {
		return eval(ast, env)
	}
}

func (eval EvalType) Trampoline(ast core.Any, env core.Env) (core.Any, error) {
	value, err := eval(ast, env)
	if err != nil {
		return nil, err
	}
	for {
		switch value.(type) {
		default:
			return value, nil
		case ThunkType:
			break
		}
		thunk := value.(ThunkType)
		next, err := thunk()
		if err != nil {
			return nil, err
		}
		value = next
	}
}
