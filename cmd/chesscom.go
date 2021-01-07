package cmd

import (
	"github.com/spf13/cobra"
	chesscom "github.com/yafred/chess-explorer/internal/chesscom"
)

var chesscomCmd = &cobra.Command{
	Use:   "chesscom [user]",
	Short: "Download games for a given user from Chess.com",
	Long:  `Download games for a given user from Chess.com`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		chesscom.DownloadGames(args[0])
	},
}

func init() {
	rootCmd.AddCommand(chesscomCmd)
}
