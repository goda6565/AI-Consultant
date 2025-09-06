package service

import (
	"errors"
	"testing"
	"time"

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

	testTitle, _ := value.NewTitle("テストドキュメント")
	testDocType, _ := value.NewDocumentType("structured")
	testDocExt, _ := value.NewDocumentExtension("pdf")
	testDateTime := sharedValue.DateTime(time.Now())
	testID := sharedValue.ID("test-id")

	existingDoc := entity.NewDocument(
		testID,
		testTitle,
		testDocType,
		testDocExt,
		nil,
		testDateTime,
	)

	tests := []struct {
		name           string
		title          value.Title
		mockSetup      func()
		expectedResult bool
		expectedError  bool
	}{
		{
			name:  "重複あり - 同じタイトルのドキュメントが存在",
			title: testTitle,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByTitle(testTitle).
					Return(existingDoc, nil).
					Times(1)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:  "重複なし - 該当するドキュメントが存在しない",
			title: testTitle,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByTitle(testTitle).
					Return(nil, nil).
					Times(1)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:  "エラー - リポジトリでエラーが発生",
			title: testTitle,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByTitle(testTitle).
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
			checker := NewDuplicateChecker(mockRepo)

			result, err := checker.CheckDuplicateByTitle(tt.title)

			if tt.expectedError && err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("エラーが期待されませんでしたが、エラーが発生しました: %v", err)
			}
			if result != tt.expectedResult {
				t.Errorf("期待された結果: %v, 実際の結果: %v", tt.expectedResult, result)
			}
		})
	}
}

func TestDuplicateChecker_CheckDuplicateByPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockDocumentRepository(ctrl)

	testPath := value.NewStoragePath("test-bucket", "test-object")
	testTitle, _ := value.NewTitle("テストドキュメント")
	testDocType, _ := value.NewDocumentType("unstructured")
	testDocExt, _ := value.NewDocumentExtension("pdf")
	testDateTime := sharedValue.DateTime(time.Now())
	testID := sharedValue.ID("test-id")

	existingDoc := entity.NewDocument(
		testID,
		testTitle,
		testDocType,
		testDocExt,
		nil,
		testDateTime,
	)

	tests := []struct {
		name           string
		path           *value.StoragePath
		mockSetup      func()
		expectedResult bool
		expectedError  bool
	}{
		{
			name: "重複あり - 同じパスのドキュメントが存在",
			path: &testPath,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByPath(testPath).
					Return(existingDoc, nil).
					Times(1)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "重複なし - 該当するドキュメントが存在しない",
			path: &testPath,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByPath(testPath).
					Return(nil, nil).
					Times(1)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name: "パスがnil - 重複なしとして扱う",
			path: nil,
			mockSetup: func() {
				// nilの場合はリポジトリは呼ばれない
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name: "エラー - リポジトリでエラーが発生",
			path: &testPath,
			mockSetup: func() {
				mockRepo.EXPECT().
					FindByPath(testPath).
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
			checker := NewDuplicateChecker(mockRepo)

			result, err := checker.CheckDuplicateByPath(tt.path)

			if tt.expectedError && err == nil {
				t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("エラーが期待されませんでしたが、エラーが発生しました: %v", err)
			}
			if result != tt.expectedResult {
				t.Errorf("期待された結果: %v, 実際の結果: %v", tt.expectedResult, result)
			}
		})
	}
}
