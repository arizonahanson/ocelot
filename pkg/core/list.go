package core

import "fmt"

type List []Any

func (value List) String() string {
	result := "("
	for i, item := range value {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%v", item)
	}
	return result + ")"
}
