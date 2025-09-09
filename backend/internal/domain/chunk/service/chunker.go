package service

import (
	"context"
	"strings"
)

// Chunk content: 50 chars before + 200 chars center + 50 chars after = max 300 chars
// Parent context: 400 chars before + 200 chars center + 400 chars after = max 1000 chars
const MaxChunkOverlap = 50
const MaxChunkSize = 200
const MaxParentChunkSize = 1000

type ChunkerInput struct {
	Text string
}

type Chunk struct {
	Content       string
	ParentContent string
}

type ChunkerOutput struct {
	Chunks []Chunk
}

type Chunker struct {
}

func NewChunkService() *Chunker {
	return &Chunker{}
}

func (c *Chunker) Execute(ctx context.Context, input ChunkerInput) (*ChunkerOutput, error) {
	if input.Text == "" {
		return &ChunkerOutput{Chunks: []Chunk{}}, nil
	}

	text := strings.TrimSpace(normalizeWhitespace(input.Text))
	if text == "" {
		return &ChunkerOutput{Chunks: []Chunk{}}, nil
	}

	runes := []rune(text)
	total := len(runes)
	if total == 0 {
		return &ChunkerOutput{Chunks: []Chunk{}}, nil
	}

	// parent = 400 + 200 + 400 when 1000 and 200
	halfParent := (MaxParentChunkSize - MaxChunkSize) / 2

	var chunks []Chunk
	for centerStart := 0; centerStart < total; centerStart += MaxChunkSize {
		centerEnd := centerStart + MaxChunkSize
		if centerEnd > total {
			centerEnd = total
		}

		contentStart := centerStart - MaxChunkOverlap
		if contentStart < 0 {
			contentStart = 0
		}
		contentEnd := centerEnd + MaxChunkOverlap
		if contentEnd > total {
			contentEnd = total
		}

		parentStart := centerStart - halfParent
		if parentStart < 0 {
			parentStart = 0
		}
		parentEnd := centerEnd + halfParent
		if parentEnd > total {
			parentEnd = total
		}

		chunks = append(chunks, Chunk{
			Content:       string(runes[contentStart:contentEnd]),
			ParentContent: string(runes[parentStart:parentEnd]),
		})
	}

	return &ChunkerOutput{Chunks: chunks}, nil
}

func normalizeWhitespace(text string) string {
	normalized := strings.ReplaceAll(text, "\n", " ")
	return strings.Join(strings.Fields(normalized), " ")
}
