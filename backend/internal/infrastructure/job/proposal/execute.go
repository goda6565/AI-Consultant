package proposal

import (
	"context"
	"os"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/errors"
	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/job"
	"github.com/goda6565/ai-consultant/backend/internal/usecase/proposal"
)

type ExecuteProposalJobApplication struct {
	executeProposalUseCase proposal.ExecuteProposalInputPort
}

func NewExecuteProposal(ctx context.Context, executeProposalUseCase proposal.ExecuteProposalInputPort) job.JobApplication {
	return &ExecuteProposalJobApplication{
		executeProposalUseCase: executeProposalUseCase,
	}
}

func (j *ExecuteProposalJobApplication) Execute(ctx context.Context) error {
	problemID := os.Getenv("PROBLEM_ID")
	if problemID == "" {
		return errors.NewInfrastructureError(errors.BadRequestError, "problem id is not set")
	}
	return j.executeProposalUseCase.Execute(ctx, proposal.ExecuteProposalUseCaseInput{
		ProblemID: problemID,
	})
}
