package cmd

import (
	"github.com/spf13/cobra"
	pgntodb "github.com/yafred/chess-explorer/internal/pgntodb"
)

var pgnToDbCmd = &cobra.Command{
	Use:   "pgntodb [pgn file]",
	Short: "Parse a pgn file and feed mongo database",
	Long:  `Parse a pgn file and feed mongo database. Designed for chess.com and lichess.org`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		latestGame := pgntodb.Game{}
		pgntodb.Process(args[0], latestGame)
	},
}

func init() {
	rootCmd.AddCommand(pgnToDbCmd)
}
