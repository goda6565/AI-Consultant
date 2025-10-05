package value

type RetryCount int

func (r RetryCount) Equals(other RetryCount) bool {
	return r == other
}

func (r RetryCount) Value() int {
	return int(r)
}

func NewRetryCount(value int) RetryCount {
	return RetryCount(value)
}
