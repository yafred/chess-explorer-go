package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yafred/chess-explorer/internal/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display build information",
	Long:  `Display build information`,
	Run: func(cmd *cobra.Command, args []string) {
		version.DisplayVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
