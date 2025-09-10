package firebase

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/pkg/auth"
)

type FirebaseClient struct {
}

func NewFirebaseClient() auth.Authenticator {
	return &FirebaseClient{}
}

func (c *FirebaseClient) Validate(ctx context.Context, token string) (*auth.AuthenticatedValue, error) {
	return nil, nil
}
