package value

type Content struct {
	value string
}

func (c *Content) Value() string {
	return c.value
}

func (c *Content) Equals(other Content) bool {
	return c.value == other.value
}

func NewContent(value string) *Content {
	return &Content{value: value}
}
