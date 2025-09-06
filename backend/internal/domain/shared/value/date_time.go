package value

import (
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

type DateTime time.Time

func (dt DateTime) Equals(other DateTime) bool {
	return time.Time(dt).Equal(time.Time(other))
}

func (dt DateTime) Value() time.Time {
	return time.Time(dt)
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(time.RFC3339)
}

func (dt DateTime) StringDate() string {
	return time.Time(dt).Format("2006-01-02")
}

func (dt DateTime) StringTime() string {
	return time.Time(dt).Format("15:04:05")
}

func (dt DateTime) StringDateTime() string {
	return time.Time(dt).Format("2006-01-02 15:04:05")
}

func (dt DateTime) Before(other DateTime) bool {
	return time.Time(dt).Before(time.Time(other))
}

func (dt DateTime) After(other DateTime) bool {
	return time.Time(dt).After(time.Time(other))
}

func NewDateTime(value string) (DateTime, error) {
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return DateTime{}, errors.NewDomainError(errors.ValidationError, "invalid datetime format, expected RFC3339 (2006-01-02T15:04:05Z07:00)")
	}

	return DateTime(parsedTime), nil
}
