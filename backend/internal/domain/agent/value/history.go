package value

type History struct {
	value string
}

func NewHistory(value string) *History {
	return &History{value: value}
}

func (h *History) GetValue() string {
	return h.value
}
