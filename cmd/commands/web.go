package commands

import (
	"github.com/akbariandev/jumpy/server/http"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run Web UI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		http.Run()
	},
}
