package value

type ActionInput struct {
	value string
}

func NewActionInput(value string) (*ActionInput, error) {
	return &ActionInput{value: value}, nil
}

func (a *ActionInput) Value() string {
	return a.value
}

func (a *ActionInput) Equals(other ActionInput) bool {
	return a.value == other.value
}
