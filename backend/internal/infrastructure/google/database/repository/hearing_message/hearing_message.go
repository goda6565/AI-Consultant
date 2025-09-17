package hearingmessage

import (
	"context"
	"fmt"

	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/hearing_message/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/google/database/internal/gen/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type HearingMessageRepository struct {
	tx   pgx.Tx
	pool *database.AppPool
}

func NewHearingMessageRepository(pool *database.AppPool) repository.HearingMessageRepository {
	return &HearingMessageRepository{tx: nil, pool: pool}
}

func (r *HearingMessageRepository) WithTx(tx pgx.Tx) *HearingMessageRepository {
	return &HearingMessageRepository{tx: tx, pool: r.pool}
}

func (r *HearingMessageRepository) FindByHearingID(ctx context.Context, hearingID sharedValue.ID) ([]entity.HearingMessage, error) {
	q := app.New(r.pool)
	var hID pgtype.UUID
	if err := hID.Scan(hearingID.Value()); err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan hearing id: %v", err))
	}
	hearingMessages, err := q.GetHearingMessageByHearingID(ctx, hID)
	if err != nil {
		return nil, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to get hearing messages by hearing id: %v", err))
	}
	entities := make([]entity.HearingMessage, len(hearingMessages))
	for i, hearingMessage := range hearingMessages {
		entity, err := toEntity(hearingMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to convert hearing message to entity: %v", err)
		}
		entities[i] = *entity
	}
	return entities, nil
}

func (r *HearingMessageRepository) Create(ctx context.Context, hearingMessage *entity.HearingMessage) error {
	q := app.New(r.pool)
	var id pgtype.UUID
	if err := id.Scan(hearingMessage.GetID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan id: %v", err))
	}
	var hearingID pgtype.UUID
	if err := hearingID.Scan(hearingMessage.GetHearingID().Value()); err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan hearing id: %v", err))
	}
	hearingMessageValue := hearingMessage.GetMessage()
	err := q.CreateHearingMessage(ctx, app.CreateHearingMessageParams{
		ID:        id,
		HearingID: hearingID,
		Role:      hearingMessage.GetRole().Value(),
		Message:   hearingMessageValue.Value(),
	})
	if err != nil {
		return errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to create hearing message: %v", err))
	}
	return nil
}

func (r *HearingMessageRepository) DeleteByHearingID(ctx context.Context, hearingID sharedValue.ID) (numDeleted int64, err error) {
	q := app.New(r.pool)
	var hID pgtype.UUID
	if err := hID.Scan(hearingID.Value()); err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to scan hearing id: %v", err))
	}
	numDeleted, err = q.DeleteHearingMessageByHearingID(ctx, hID)
	if err != nil {
		return 0, errors.NewInfrastructureError(errors.ExternalServiceError, fmt.Sprintf("failed to delete hearing messages by hearing id: %v", err))
	}
	return numDeleted, nil
}

func toEntity(hearingMessage app.HearingMessage) (*entity.HearingMessage, error) {
	id, err := sharedValue.NewID(hearingMessage.ID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	hearingID, err := sharedValue.NewID(hearingMessage.HearingID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create hearing id: %w", err)
	}
	role, err := value.NewRole(hearingMessage.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	message, err := value.NewMessage(hearingMessage.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	createdAt := hearingMessage.CreatedAt.Time

	return entity.NewHearingMessage(id, hearingID, role, *message, &createdAt), nil
}
