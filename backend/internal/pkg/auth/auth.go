package auth

import "context"

type AuthenticatedValue struct {
	Token  string
	UserID string
	Email  string
}

type Authenticator interface {
	Validate(ctx context.Context, token string) (*AuthenticatedValue, error)
}
