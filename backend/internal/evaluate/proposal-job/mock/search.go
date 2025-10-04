package mock

import (
	"github.com/goda6565/ai-consultant/backend/internal/domain/search"
	searchMock "github.com/goda6565/ai-consultant/backend/internal/domain/search/mock"
	gomock "go.uber.org/mock/gomock"
)

func NewMockDocumentSearchClient() (search.DocumentSearchClient, func()) {
	ctrl := gomock.NewController(nil)
	m := searchMock.NewMockDocumentSearchClient(ctrl)
	m.EXPECT().Search(gomock.Any(), gomock.Any()).Return(&search.DocumentSearchOutput{}, nil).AnyTimes()
	return m, func() {
		ctrl.Finish()
	}
}
