package service

import (
	"context"
	"errors"
	"testing"

	"github.com/goda6565/ai-consultant/backend/internal/domain/document/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/repository/mock"
	"github.com/goda6565/ai-consultant/backend/internal/domain/document/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"go.uber.org/mock/gomock"
)

func TestDuplicateChecker_CheckDuplicateByTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockDocumentRepository(ctrl)

	testTitle, _ := value.NewTitle("Test Document")
	testDocExt, _ := value.NewDocumentType("pdf")
	testID := sharedValue.ID("test-id")
	testStoragePath := value.NewStorageInfo("test-bucket", "test-object")

	existingDoc := entity.NewDocument(
		testID,
		testTitle,
		testDocExt,
		testStoragePath,
		value.DocumentStatusProcessing,
		value.NewRetryCount(0),
		nil,
		nil,
	)

	tests := []struct {
		name           string
		title          value.Title
		mockSetup      func()
		expectedResult bool
		expectedError  bool
	}{
		{
			name:  "duplicate found - document with same title exists",
			title: testTitle,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByTitle(gomock.Any(), testTitle).
					Return(existingDoc, nil).
					Times(1)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:  "no duplicate - no matching document exists",
			title: testTitle,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByTitle(gomock.Any(), testTitle).
					Return(nil, nil).
					Times(1)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:  "error - repository error occurred",
			title: testTitle,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByTitle(gomock.Any(), testTitle).
					Return(nil, errors.New("database error")).
					Times(1)
			},
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			checker := NewDuplicateCheckService(mockRepo)

			result, err := checker.CheckDuplicateByTitle(context.Background(), tt.title)

			if tt.expectedError && err == nil {
				t.Errorf("expected error but no error occurred")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error occurred: %v", err)
			}
			if result != tt.expectedResult {
				t.Errorf("expected result: %v, actual result: %v", tt.expectedResult, result)
			}
		})
	}
}
