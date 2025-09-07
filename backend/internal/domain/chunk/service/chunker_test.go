package service

import (
	"context"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestChunker_EmptyText(t *testing.T) {
	// Empty text should return empty chunks
	ch := NewChunker()
	out, err := ch.Execute(context.Background(), ChunkerInput{Text: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Chunks) != 0 {
		t.Fatalf("expected 0 chunks, got %d", len(out.Chunks))
	}
}

func TestChunker_FixedWindow_WithUniformRunesSmall(t *testing.T) {
	// 600 runes: centers at 0, 200, 400
	text := strings.Repeat("あ", 600)
	ch := NewChunker()
	out, err := ch.Execute(context.Background(), ChunkerInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantChunks := (utf8.RuneCountInString(text) + MaxChunkSize - 1) / MaxChunkSize
	if len(out.Chunks) != wantChunks {
		t.Fatalf("chunk count mismatch: got %d want %d", len(out.Chunks), wantChunks)
	}

	// Expected content lengths: [250, 300, 250]
	lens := []int{utf8.RuneCountInString(out.Chunks[0].Content), utf8.RuneCountInString(out.Chunks[1].Content), utf8.RuneCountInString(out.Chunks[2].Content)}
	if lens[0] != 250 || lens[1] != 300 || lens[2] != 250 {
		t.Errorf("unexpected content lengths: %v", lens)
	}

	// Expected parent contents
	expectedParentContents := []string{strings.Repeat("あ", 600), strings.Repeat("あ", 600), strings.Repeat("あ", 600)}
	for i, c := range out.Chunks {
		if c.ParentContent != expectedParentContents[i] {
			t.Errorf("parent content mismatch at %d: got %q want %q", i, c.ParentContent, expectedParentContents[i])
		}
	}
}

func TestChunker_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		wantChunks   int
		validateFunc func(t *testing.T, chunks []Chunk)
	}{
		{
			name:       "exactly_200_chars",
			text:       strings.Repeat("A", 200),
			wantChunks: 1,
			validateFunc: func(t *testing.T, chunks []Chunk) {
				// Only one chunk, content and parent should be identical
				if chunks[0].Content != chunks[0].ParentContent {
					t.Errorf("single chunk: content != parent")
				}
				if utf8.RuneCountInString(chunks[0].Content) != 200 {
					t.Errorf("single chunk length should be 200, got %d", utf8.RuneCountInString(chunks[0].Content))
				}
			},
		},
		{
			name:       "exactly_250_chars",
			text:       strings.Repeat("B", 250),
			wantChunks: 2,
			validateFunc: func(t *testing.T, chunks []Chunk) {
				// First chunk: 0-250 (no left overlap)
				if utf8.RuneCountInString(chunks[0].Content) != 250 {
					t.Errorf("first chunk should be 250, got %d", utf8.RuneCountInString(chunks[0].Content))
				}
				// Second chunk: 150-250 (50 chars)
				if utf8.RuneCountInString(chunks[1].Content) != 100 {
					t.Errorf("second chunk should be 100, got %d", utf8.RuneCountInString(chunks[1].Content))
				}
			},
		},
		{
			name:       "whitespace_normalization",
			text:       "Line1\n\nLine2\t\tLine3   Line4",
			wantChunks: 1,
			validateFunc: func(t *testing.T, chunks []Chunk) {
				expected := "Line1 Line2 Line3 Line4"
				if chunks[0].Content != expected {
					t.Errorf("whitespace not normalized: got %q want %q", chunks[0].Content, expected)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := NewChunker()
			out, err := ch.Execute(context.Background(), ChunkerInput{Text: tt.text})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(out.Chunks) != tt.wantChunks {
				t.Fatalf("chunk count mismatch: got %d want %d", len(out.Chunks), tt.wantChunks)
			}
			tt.validateFunc(t, out.Chunks)
		})
	}
}

func TestChunker_OverlapValidation(t *testing.T) {
	// Test with 450 chars: centers at 0, 200, 400
	text := strings.Repeat("X", 450)
	ch := NewChunker()
	out, err := ch.Execute(context.Background(), ChunkerInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(out.Chunks))
	}

	// Expected content lengths
	expectedLengths := []int{250, 300, 100}
	for i, chunk := range out.Chunks {
		actualLen := utf8.RuneCountInString(chunk.Content)
		if actualLen != expectedLengths[i] {
			t.Errorf("chunk %d length mismatch: got %d want %d", i, actualLen, expectedLengths[i])
		}
	}

	// Expected parent contents
	expectedParentContents := []int{450, 450, 450}
	for i, c := range out.Chunks {
		actualParentLen := utf8.RuneCountInString(c.ParentContent)
		if actualParentLen != expectedParentContents[i] {
			t.Errorf("parent content length mismatch at %d: got %d want %d", i, actualParentLen, expectedParentContents[i])
		}
	}
}

func TestChunker_OverlapValidation_LongText(t *testing.T) {
	// Test with 1500 chars: centers at 0, 200, 400, 600, 800, 1000, 1200, 1400
	text := strings.Repeat("X", 1500)
	ch := NewChunker()
	out, err := ch.Execute(context.Background(), ChunkerInput{Text: text})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Chunks) != 8 {
		t.Fatalf("expected 8 chunks, got %d", len(out.Chunks))
	}

	// Expected content lengths for each chunk:
	// Chunk 0 (center 0-200): content 0-250 = 250 chars
	// Chunk 1 (center 200-400): content 150-450 = 300 chars
	// Chunk 2 (center 400-600): content 350-650 = 300 chars
	// Chunk 3 (center 600-800): content 550-850 = 300 chars
	// Chunk 4 (center 800-1000): content 750-1050 = 300 chars
	// Chunk 5 (center 1000-1200): content 950-1250 = 300 chars
	// Chunk 6 (center 1200-1400): content 1150-1450 = 300 chars
	// Chunk 7 (center 1400-1500): content 1350-1500 = 150 chars
	expectedLengths := []int{250, 300, 300, 300, 300, 300, 300, 150}
	for i, chunk := range out.Chunks {
		actualLen := utf8.RuneCountInString(chunk.Content)
		if actualLen != expectedLengths[i] {
			t.Errorf("chunk %d length mismatch: got %d want %d", i, actualLen, expectedLengths[i])
		}
	}

	// Expected parent content lengths:
	// For MaxParentChunkSize=1000, halfParent=(1000-200)/2=400
	// Chunk 0 (center 0-200): parent 0-600 = 600 chars
	// Chunk 1 (center 200-400): parent 0-800 = 800 chars
	// Chunk 2 (center 400-600): parent 0-1000 = 1000 chars
	// Chunk 3 (center 600-800): parent 200-1200 = 1000 chars
	// Chunk 4 (center 800-1000): parent 400-1400 = 1000 chars
	// Chunk 5 (center 1000-1200): parent 600-1500 = 900 chars
	// Chunk 6 (center 1200-1400): parent 800-1500 = 700 chars
	// Chunk 7 (center 1400-1500): parent 1000-1500 = 500 chars
	expectedParentLengths := []int{600, 800, 1000, 1000, 1000, 900, 700, 500}
	for i, chunk := range out.Chunks {
		actualParentLen := utf8.RuneCountInString(chunk.ParentContent)
		if actualParentLen != expectedParentLengths[i] {
			t.Errorf("chunk %d parent length mismatch: got %d want %d", i, actualParentLen, expectedParentLengths[i])
		}
	}
}
