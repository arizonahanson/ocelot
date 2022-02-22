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
	more := slice(rest)
	result := make([]core.Any, len(more)+1)
	result[0] = first
	for i, group := range more {
		next := slice(group)[index]
		result[i+1] = next
	}
	return result
}

func merge(first, rest interface{}, keyIndex int, valueIndex int) map[core.Key]core.Any {
	more := slice(rest)
	result := make(map[core.Key]core.Any, len(more)+1)
	pair := slice(first)
	if pair == nil {
		return result
	}
	result[pair[keyIndex].(core.Key)] = pair[valueIndex]
	for _, group := range more {
		pair := slice(group)
		result[pair[keyIndex+1].(core.Key)] = pair[valueIndex+1]
	}
	return result
}

func pos(p position) *core.Position {
	return &core.Position{Line: p.line, Col: p.col, Offset: p.offset}
}
