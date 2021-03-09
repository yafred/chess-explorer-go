package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yafred/chess-explorer/internal/cache"
	chesscom "github.com/yafred/chess-explorer/internal/chesscom"
)

var chesscomPgn string

var chesscomCmd = &cobra.Command{
	Use:   "chesscom [user]",
	Short: "Download games for a given user from Chess.com",
	Long:  `Download games for a given user from Chess.com`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		chesscom.DownloadGames(args[0], chesscomPgn)
		cache.UpdateInitialValues()
	},
}

func init() {
	rootCmd.AddCommand(chesscomCmd)

	chesscomCmd.Flags().StringVar(&chesscomPgn, "keep", "", "file where the PGN will be kept")
}
