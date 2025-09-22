package value

type StatusMessage struct {
	value string
}

func NewStatusMessage(value string) *StatusMessage {
	return &StatusMessage{value: value}
}

func (s *StatusMessage) Value() string {
	return s.value
}

func (s *StatusMessage) Equals(other StatusMessage) bool {
	return s.value == other.value
}
