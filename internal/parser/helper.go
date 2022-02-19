package parser

import "github.com/starlight/ocelot/pkg/core"

// cast to []interface{}
func slice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// build []core.Any from first, rest=[[_, next], ...]
func join(first, rest interface{}, index int) []core.Any {
	if first == nil {
		return []core.Any{}
	}
	result := []core.Any{first}
	for _, group := range slice(rest) {
		next := slice(group)[index]
		result = append(result, next)
	}
	return result
}

func merge(first, rest interface{}, keyIndex int, valueIndex int) map[core.Key]core.Any {
	result := make(map[core.Key]core.Any)
	pair := slice(first)
	if pair == nil {
		return result
	}
	result[pair[keyIndex].(core.Key)] = pair[valueIndex]
	for _, group := range slice(rest) {
		pair := slice(group)
		result[pair[keyIndex+1].(core.Key)] = pair[valueIndex+1]
	}
	return result
}

func pos(p position) *core.Position {
	return &core.Position{Line: p.line, Col: p.col, Offset: p.offset}
}
