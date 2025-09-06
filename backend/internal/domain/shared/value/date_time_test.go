package value

import (
	"errors"
	"testing"
	"time"

	domainErrors "github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

func TestNewDateTime(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError bool
	}{
		{
			name:          "valid RFC3339 with Z",
			input:         "2023-12-25T10:30:00Z",
			expectedError: false,
		},
		{
			name:          "valid RFC3339 with timezone",
			input:         "2023-12-25T10:30:00+09:00",
			expectedError: false,
		},
		{
			name:          "valid RFC3339 with negative timezone",
			input:         "2023-12-25T10:30:00-05:00",
			expectedError: false,
		},
		{
			name:          "invalid format - missing timezone",
			input:         "2023-12-25T10:30:00",
			expectedError: true,
		},
		{
			name:          "invalid format - space instead of T",
			input:         "2023-12-25 10:30:00Z",
			expectedError: true,
		},
		{
			name:          "empty string",
			input:         "",
			expectedError: true,
		},
		{
			name:          "invalid date",
			input:         "2023-13-32T10:30:00Z",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := NewDateTime(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				var domainErr *domainErrors.DomainError
				if err != nil && !errors.As(err, &domainErr) {
					t.Errorf("expected domain error but got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if dt.String() == "" {
					t.Error("expected non-empty string representation")
				}
			}
		})
	}
}

func TestDateTime_Equals(t *testing.T) {
	dt1, _ := NewDateTime("2023-12-25T10:30:00Z")
	dt2, _ := NewDateTime("2023-12-25T10:30:00Z")
	dt3, _ := NewDateTime("2023-12-25T10:30:01Z")

	if !dt1.Equals(dt2) {
		t.Error("expected equal datetimes to be equal")
	}

	if dt1.Equals(dt3) {
		t.Error("expected different datetimes to not be equal")
	}
}

func TestDateTime_BeforeAfter(t *testing.T) {
	earlier, _ := NewDateTime("2023-12-25T10:30:00Z")
	later, _ := NewDateTime("2023-12-25T10:30:01Z")

	if !earlier.Before(later) {
		t.Error("expected earlier datetime to be before later datetime")
	}

	if !later.After(earlier) {
		t.Error("expected later datetime to be after earlier datetime")
	}

	if earlier.After(later) {
		t.Error("expected earlier datetime to not be after later datetime")
	}

	if later.Before(earlier) {
		t.Error("expected later datetime to not be before earlier datetime")
	}
}

func TestDateTime_Value(t *testing.T) {
	input := "2023-12-25T15:30:45Z"
	dt, _ := NewDateTime(input)

	actual := dt.Value()
	expected := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	if !actual.Equal(expected) {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestDateTime_String(t *testing.T) {
	input := "2023-12-25T15:30:45Z"
	dt, _ := NewDateTime(input)

	actual := dt.String()
	expected := input
	if actual != expected {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}

func TestDateTime_StringDate(t *testing.T) {
	dt, _ := NewDateTime("2023-12-25T15:30:45Z")

	actual := dt.StringDate()
	expected := "2023-12-25"
	if actual != expected {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}

func TestDateTime_StringTime(t *testing.T) {
	dt, _ := NewDateTime("2023-12-25T15:30:45Z")

	actual := dt.StringTime()
	expected := "15:30:45"
	if actual != expected {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}

func TestDateTime_StringDateTime(t *testing.T) {
	dt, _ := NewDateTime("2023-12-25T15:30:45Z")

	actual := dt.StringDateTime()
	expected := "2023-12-25 15:30:45"
	if actual != expected {
		t.Errorf("expected %s but got %s", expected, actual)
	}
}

func TestDateTime_TimezoneHandling(t *testing.T) {
	// 同じ時刻の異なるタイムゾーン表現
	utc, _ := NewDateTime("2023-12-25T15:30:00Z")
	jst, _ := NewDateTime("2023-12-26T00:30:00+09:00")

	if !utc.Equals(jst) {
		t.Error("expected UTC and JST times representing the same moment to be equal")
	}
}
