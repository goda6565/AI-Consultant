package value

type Answered struct {
	value bool
}

func NewAnswered(value bool) *Answered {
	return &Answered{value: value}
}

func (a *Answered) Equals(other *Answered) bool {
	return a.value == other.value
}

func (a *Answered) Value() bool {
	return a.value
}
