package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yafred/chess-explorer/internal/try"
)

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "Try something",
	Long:  `Try something`,
	Run: func(cmd *cobra.Command, args []string) {
		try.Something()
	},
}

func init() {
	rootCmd.AddCommand(tryCmd)
}
