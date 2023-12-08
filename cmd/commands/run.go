package commands

import (
	"github.com/akbariandev/jumpy/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a DAG Node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		app.Start(listenPort, hostGroupName)
	},
}
