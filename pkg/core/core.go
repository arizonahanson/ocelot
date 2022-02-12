package core

type Any interface{}

type Symbol string

type Vector []Any

type Environment map[Symbol]Any

type Function func(args ...Any) Any
