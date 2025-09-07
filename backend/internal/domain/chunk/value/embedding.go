package value

import (
	"reflect"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

const VectorLength = 1536

type Embedding []float32

func (e Embedding) Equals(other Embedding) bool {
	return reflect.DeepEqual(e, other)
}

func (e Embedding) Value() []float32 {
	return e
}

func NewEmbedding(value []float32) (Embedding, error) {
	if len(value) != VectorLength {
		return nil, errors.NewDomainError(errors.ValidationError, "embedding must be 1536 elements long")
	}
	return Embedding(value), nil
}
