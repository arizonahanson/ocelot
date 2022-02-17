package base

import "github.com/starlight/ocelot/pkg/core"

type Function func(args []core.Any, env Env) (core.Any, error)

type Thunk func() (core.Any, error)

func (fn Function) thunk(args []core.Any, env Env) Thunk {
	return func() (core.Any, error) {
		return fn(args, env)
	}
}

func (fn Function) trampoline() Function {
	return func(args []core.Any, env Env) (core.Any, error) {
		res, err := fn(args, env)
		if err != nil {
			return nil, err
		}
		for {
			switch res.(type) {
			default:
				return res, nil
			case Thunk:
				fn2 := res.(Thunk)
				res2, err := fn2()
				if err != nil {
					return nil, err
				}
				res = res2
			}
		}
	}
}
