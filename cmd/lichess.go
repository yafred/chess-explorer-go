package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yafred/chess-explorer/internal/lichess"
)

var lichessCmd = &cobra.Command{
	Use:   "lichess [user]",
	Short: "Download games for a given user from Lichess.org",
	Long:  `Download games for a given user from Lichess.org`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lichess.DownloadGames(args[0])
	},
}

func init() {
	rootCmd.AddCommand(lichessCmd)
}
