package problem

import (
	"context"
	"fmt"

	jobConfigEntity "github.com/goda6565/ai-consultant/backend/internal/domain/job_config/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/entity"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/repository"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/service"
	"github.com/goda6565/ai-consultant/backend/internal/domain/problem/value"
	problemFieldEntity "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/entity"
	problemFieldRepository "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/repository"
	problemFieldService "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/service"
	problemFieldValue "github.com/goda6565/ai-consultant/backend/internal/domain/problem_field/value"
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
	"github.com/goda6565/ai-consultant/backend/internal/pkg/uuid"
	transaction "github.com/goda6565/ai-consultant/backend/internal/usecase/ports/transaction"
)

type CreateProblemInputPort interface {
	Execute(ctx context.Context, input CreateProblemUseCaseInput) (*CreateProblemOutput, error)
}

type CreateProblemUseCaseInput struct {
	Description string
}

type CreateProblemOutput struct {
	Problem *entity.Problem
}

type CreateProblemInteractor struct {
	generateTitleService   *service.GenerateTitleService
	problemRepository      repository.ProblemRepository
	problemFieldRepository problemFieldRepository.ProblemFieldRepository
	problemFieldService    *problemFieldService.GenerateProblemFieldService
	adminUnitOfWork        transaction.AdminUnitOfWork
}

func NewCreateProblemUseCase(generateTitleService *service.GenerateTitleService, problemRepository repository.ProblemRepository, problemFieldRepository problemFieldRepository.ProblemFieldRepository, problemFieldService *problemFieldService.GenerateProblemFieldService, adminUnitOfWork transaction.AdminUnitOfWork) CreateProblemInputPort {
	return &CreateProblemInteractor{generateTitleService: generateTitleService, problemRepository: problemRepository, problemFieldRepository: problemFieldRepository, problemFieldService: problemFieldService, adminUnitOfWork: adminUnitOfWork}
}

func (i *CreateProblemInteractor) Execute(ctx context.Context, input CreateProblemUseCaseInput) (*CreateProblemOutput, error) {
	description, err := value.NewDescription(input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create description: %w", err)
	}

	// generate title
	title, err := i.generateTitleService.Execute(ctx, service.GenerateTitleServiceInput{Description: input.Description})
	if err != nil {
		return nil, fmt.Errorf("failed to generate title: %w", err)
	}

	// create value objects
	id, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return nil, fmt.Errorf("failed to create id: %w", err)
	}
	titleValue, err := value.NewTitle(title.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create title: %w", err)
	}

	// create problem
	problem := entity.NewProblem(id, *titleValue, *description, value.StatusPending, nil)

	// generate problem fields
	problemFields, err := i.problemFieldService.Execute(ctx, problemFieldService.GenerateProblemFieldServiceInput{Problem: problem})
	if err != nil {
		return nil, fmt.Errorf("failed to generate problem fields: %w", err)
	}

	// create problem fields
	problemFieldsEntities := make([]*problemFieldEntity.ProblemField, len(problemFields.Fields))
	for i, problemField := range problemFields.Fields {
		id, err := sharedValue.NewID(uuid.NewUUID())
		if err != nil {
			return nil, fmt.Errorf("failed to create id: %w", err)
		}
		field, err := problemFieldValue.NewField(problemField)
		if err != nil {
			return nil, fmt.Errorf("failed to create field: %w", err)
		}
		problemFieldsEntities[i] = problemFieldEntity.NewProblemField(id, problem.GetID(), *field, *problemFieldValue.NewAnswered(false), nil)
	}

	// create job config
	jobConfigID, err := sharedValue.NewID(uuid.NewUUID())
	if err != nil {
		return nil, fmt.Errorf("failed to create job config id: %w", err)
	}
	jobConfig := jobConfigEntity.NewJobConfig(jobConfigID, problem.GetID(), false)

	// save problem and problem fields in transaction
	err = i.adminUnitOfWork.WithTx(ctx, func(ctx context.Context) error {
		// save problem first (parent)
		if err := i.adminUnitOfWork.ProblemRepository(ctx).Create(ctx, problem); err != nil {
			return fmt.Errorf("failed to save problem: %w", err)
		}

		// then save problem fields (children)
		for _, problemField := range problemFieldsEntities {
			if err := i.adminUnitOfWork.ProblemFieldRepository(ctx).Create(ctx, problemField); err != nil {
				return fmt.Errorf("failed to save problem field: %w", err)
			}
		}

		// then save job config
		if err := i.adminUnitOfWork.JobConfigRepository(ctx).Create(ctx, jobConfig); err != nil {
			return fmt.Errorf("failed to save job config: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save problem and problem fields: %w", err)
	}

	return &CreateProblemOutput{Problem: problem}, nil
}
