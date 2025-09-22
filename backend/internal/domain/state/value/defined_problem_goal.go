package value

type DefinedProblemGoal struct {
	value string
}

func NewDefinedProblemGoal(value string) (*DefinedProblemGoal, error) {
	return &DefinedProblemGoal{value: value}, nil
}

func (d *DefinedProblemGoal) Value() string {
	return d.value
}
