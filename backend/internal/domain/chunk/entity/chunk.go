package entity

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/chunk/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type Chunk struct {
	id            sharedValue.ID
	documentID    sharedValue.ID
	content       value.Content
	parentContent value.Content
	embedding     value.Embedding
}

func (c *Chunk) GetID() sharedValue.ID {
	return c.id
}

func (c *Chunk) GetDocumentID() sharedValue.ID {
	return c.documentID
}

func (c *Chunk) GetContent() value.Content {
	return c.content
}

func (c *Chunk) GetParentContent() value.Content {
	return c.parentContent
}

func (c *Chunk) GetEmbedding() value.Embedding {
	return c.embedding
}

func NewChunk(id sharedValue.ID, documentID sharedValue.ID, content value.Content, parentContent value.Content, embedding value.Embedding) *Chunk {
	return &Chunk{id: id, documentID: documentID, content: content, parentContent: parentContent, embedding: embedding}
}
