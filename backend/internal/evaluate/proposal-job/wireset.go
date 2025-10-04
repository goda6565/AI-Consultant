package proposaljob

import (
	llmasjudge "github.com/goda6565/ai-consultant/backend/internal/evaluate/proposal-job/llm-as-a-judge"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewProposalJobEval,
	llmasjudge.NewJudge,
)
