package firebase

import (
	"context"
	"strings"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/auth"
)

type FirebaseClient struct {
	client *firebaseauth.Client
}

func NewFirebaseClient(ctx context.Context, e *environment.Environment) auth.Authenticator {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: e.ProjectID,
	})
	if err != nil {
		panic(err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		panic(err)
	}
	return &FirebaseClient{client: client}
}

func (c *FirebaseClient) Validate(ctx context.Context, token string) (*auth.AuthenticatedValue, error) {
	// Bearer <token>
	trimmedToken := strings.TrimPrefix(token, "Bearer ")

	// Verify ID token
	idToken, err := c.client.VerifyIDToken(ctx, trimmedToken)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ForbiddenError, "invalid token")
	}

	userID := idToken.UID
	var email string
	if v, ok := idToken.Claims["email"]; ok {
		if s, ok := v.(string); ok {
			email = s
		}
	}

	return &auth.AuthenticatedValue{
		Token:  token,
		UserID: userID,
		Email:  email,
	}, nil
}
