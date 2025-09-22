package value

type DefinedProblemSummary struct {
	value string
}

func NewDefinedProblemSummary(value string) (*DefinedProblemSummary, error) {
	return &DefinedProblemSummary{value: value}, nil
}

func (d *DefinedProblemSummary) Value() string {
	return d.value
}

func (d *DefinedProblemSummary) Equals(other *DefinedProblemSummary) bool {
	return d.value == other.value
}
