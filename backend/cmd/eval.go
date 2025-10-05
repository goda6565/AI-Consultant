package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/goda6565/ai-consultant/backend/di"
)

func newProposalJobEvalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-job-eval",
		Short: "Proposal Job Eval",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run the proposal job eval",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			eval, cleanup, err := di.InitProposalJobEval(ctx)
			if err != nil {
				panic(err)
			}
			defer cleanup()
			eval.Evaluate(ctx)
		},
	})
	return cmd
}
