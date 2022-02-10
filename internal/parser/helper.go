package parser

import "github.com/starlight/ocelot/pkg/core"

// cast to []interface{}
func slice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// build []interface{} from first, rest=[[_, next], ...]
func join(first, rest interface{}, index int) []core.Any {
	result := []core.Any{first}
	for _, group := range slice(rest) {
		next := slice(group)[index]
		result = append(result, next)
	}
	return result
}
