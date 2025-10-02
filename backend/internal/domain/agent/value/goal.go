package value

type Goal struct {
	value string
}

func (c *Goal) Value() string {
	return c.value
}

func (c *Goal) Equals(other Goal) bool {
	return c.value == other.value
}

func NewGoal(value string) *Goal {
	return &Goal{value: value}
}
