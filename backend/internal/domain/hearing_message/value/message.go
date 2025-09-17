package value

type Message struct {
	value string
}

func (m *Message) Equals(other Message) bool {
	return m.value == other.value
}

func (m *Message) Value() string {
	return m.value
}

func NewMessage(value string) (*Message, error) {
	return &Message{value: value}, nil
}
