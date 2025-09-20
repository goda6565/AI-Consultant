package value

import "fmt"

type Field struct {
	value string
}

func NewField(value string) (*Field, error) {
	if value == "" {
		return nil, fmt.Errorf("value is required")
	}
	return &Field{value: value}, nil
}

func (f *Field) Value() string {
	return f.value
}

func (f *Field) Equals(other *Field) bool {
	return f.value == other.value
}
