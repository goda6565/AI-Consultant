package cmd

// import (
// 	"context"

// 	"github.com/goda6565/ai-consultant/backend/di"
// 	"github.com/spf13/cobra"
// )

// func newAgentCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "agent",
// 		Short: "Agent",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			cmd.HelpFunc()(cmd, args)
// 		},
// 	}
// 	cmd.AddCommand(&cobra.Command{
// 		Use:   "run",
// 		Short: "Run the agent server",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			ctx := context.Background()
// 			app, cleanup, err := di.InitAgentApplication(ctx)
// 			if err != nil {
// 				panic(err)
// 			}
// 			defer cleanup()
// 			app.StartApp()
// 			app.StopApp(ctx)
// 		},
// 	})
// 	return cmd
// }
