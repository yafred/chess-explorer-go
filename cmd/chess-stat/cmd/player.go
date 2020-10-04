package cmd

import (
	"github.com/spf13/cobra"
	stat "github.com/yafred/chess-com/internal/chess-stat"
)

var cachePath string
var cacheRefresh bool

var playerCmd = &cobra.Command{
	Use:   "player [username]",
	Short: "Creates stats for a chess.com player",
	Long:  `Creates stats for a chess.com player based on games downloaded from https://chess.com/`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stat.StatsToConsole(args[0], cachePath, cacheRefresh)
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)

	playerCmd.Flags().StringVarP(&cachePath, "cache-path", "c", "", "Folder where downloaded data should be kept (data will not be kept if flag absent)")
	playerCmd.MarkFlagDirname("cache-path")
	playerCmd.Flags().BoolVarP(&cacheRefresh, "cache-refresh", "r", false, "Refresh cache before executing command (if flag absent, existing data will be used)")
}
