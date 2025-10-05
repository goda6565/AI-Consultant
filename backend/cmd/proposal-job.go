package cmd

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/di"
	"github.com/spf13/cobra"
)

func newProposalJobCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-job",
		Short: "Proposal Job",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run the proposal job",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			job, cleanup, err := di.InitProposalJob(ctx)
			if err != nil {
				panic(err)
			}
			defer cleanup()
			job.Run(ctx)
		},
	})
	return cmd
}
