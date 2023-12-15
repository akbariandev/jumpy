package commands

import (
	"github.com/spf13/cobra"
	"os"
)

var nodes int
var hostGroupName string

var rootCmd = &cobra.Command{
	Use:   "ChainOfBots",
	Short: "",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&nodes, "numOfNodes", "n", 2, "number of nodes to run")
	rootCmd.PersistentFlags().StringVarP(&hostGroupName, "hostGroupname", "g", "", "name of host group")
}
