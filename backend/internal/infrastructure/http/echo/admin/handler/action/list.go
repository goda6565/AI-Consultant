package action

import (
	"context"
	"time"

	"github.com/goda6565/ai-consultant/backend/internal/domain/action/entity"
	gen "github.com/goda6565/ai-consultant/backend/internal/infrastructure/http/echo/admin/internal"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/action"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type ListActionHandler struct {
	listActionInputPort action.ListActionInputPort
}

func NewListActionHandler(listActionInputPort action.ListActionInputPort) *ListActionHandler {
	return &ListActionHandler{listActionInputPort: listActionInputPort}
}

func (h *ListActionHandler) ListActions(ctx context.Context, request gen.ListActionsRequestObject) (gen.ListActionsResponseObject, error) {
	actions, err := h.listActionInputPort.Execute(ctx, action.ListActionUseCaseInput{
		ProblemID: request.ProblemId.String(),
	})
	if err != nil {
		return nil, err
	}
	return gen.ListActions200JSONResponse{
		ListActionsSuccessJSONResponse: gen.ListActionsSuccessJSONResponse{
			Actions: toActionsJSONResponse(actions.Actions),
		},
	}, nil
}

func toActionsJSONResponse(actions []entity.Action) []gen.Action {
	actionsJSON := make([]gen.Action, len(actions))
	for i, action := range actions {
		actionsJSON[i] = toSingleActionJSON(&action)
	}
	return actionsJSON
}

func toSingleActionJSON(action *entity.Action) gen.Action {
	id := openapi_types.UUID(uuid.MustParse(action.GetID().Value()))
	problemID := openapi_types.UUID(uuid.MustParse(action.GetProblemID().Value()))
	actionType := gen.ActionType(action.GetActionType().Value())

	var createdAt time.Time
	if action.GetCreatedAt() != nil {
		createdAt = *action.GetCreatedAt()
	}

	input := action.GetInput()
	output := action.GetOutput()

	return gen.Action{
		Id:         id,
		ProblemId:  problemID,
		ActionType: actionType,
		Input:      input.Value(),
		Output:     output.Value(),
		CreatedAt:  createdAt,
	}
}
