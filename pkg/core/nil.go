package core

type Nil struct{}

func (value Nil) String() string {
	return "nil"
}
