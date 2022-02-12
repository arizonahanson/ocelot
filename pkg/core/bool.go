package core

type Bool bool

func (value Bool) String() string {
	if value {
		return "true"
	} else {
		return "false"
	}
}
