package value

type Description struct {
	value string
}

func (d Description) Equals(other Description) bool {
	return d.value == other.value
}

func (d Description) Value() string {
	return d.value
}

func NewDescription(value string) (*Description, error) {
	return &Description{value: value}, nil
}
