package value

type DefinedProblemHearingInformation struct {
	value string
}

func NewDefinedProblemHearingInformation(value string) (*DefinedProblemHearingInformation, error) {
	return &DefinedProblemHearingInformation{value: value}, nil
}

func (d *DefinedProblemHearingInformation) Value() string {
	return d.value
}
