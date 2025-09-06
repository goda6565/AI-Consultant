package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestDomainError_As(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "DomainError should be detected as DomainError",
			err:      NewDomainError(ValidationError, "test message"),
			expected: true,
		},
		{
			name:     "wrapped DomainError should be detected as DomainError",
			err:      fmt.Errorf("wrapped: %w", NewDomainError(ValidationError, "test message")),
			expected: true,
		},
		{
			name:     "standard error should not be detected as DomainError",
			err:      errors.New("standard error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var domainErr *DomainError
			result := errors.As(tt.err, &domainErr)

			if result != tt.expected {
				t.Errorf("errors.As() = %v, expected %v", result, tt.expected)
			}

			if tt.expected && domainErr == nil {
				t.Error("expected domainErr to be non-nil when errors.As returns true")
			}

			if !tt.expected && domainErr != nil {
				t.Error("expected domainErr to be nil when errors.As returns false")
			}
		})
	}
}
