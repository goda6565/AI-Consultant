package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "App",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	cmd.AddCommand(newAdminCommand())
	cmd.AddCommand(newVectorCommand())
	cmd.AddCommand(newAgentCommand())
	cmd.AddCommand(newProposalJobCommand())
	cmd.AddCommand(newProposalJobEvalCommand())
	return cmd
}
