package value

import "github.com/goda6565/ai-consultant/backend/internal/domain/errors"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

func (r Role) Equals(other Role) bool {
	return r == other
}

func (r Role) Value() string {
	return string(r)
}

func NewRole(value string) (Role, error) {
	switch value {
	case "user":
		return RoleUser, nil
	case "assistant":
		return RoleAssistant, nil
	default:
		return "", errors.NewDomainError(errors.ValidationError, "invalid role")
	}
}
