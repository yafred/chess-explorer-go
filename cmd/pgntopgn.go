package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yafred/chess-explorer/internal/pgntopgn"
)

var pgnToPgnCmd = &cobra.Command{
	Use:   "pgntopgn [pgn file]",
	Short: "Filter a pgn file",
	Long:  `Filter a pgn file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pgntopgn.Process(args[0])
	},
}

func init() {
	rootCmd.AddCommand(pgnToPgnCmd)
}
