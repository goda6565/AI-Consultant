package value

type Content string

func (c Content) Equals(other Content) bool {
	return c == other
}

func (c Content) Value() string {
	return string(c)
}

func NewContent(value string) (Content, error) {
	return Content(value), nil
}
