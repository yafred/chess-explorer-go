package cmd

import (
	"github.com/spf13/cobra"
	stat "github.com/yafred/chess-com/internal/chess-stat"
)

var playerCmd = &cobra.Command{
	Use:   "player [username]",
	Short: "Gives stats on a chess.com player",
	Long:  `Gives stats on a chess.com player based on games downloaded from https://chess.com/`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stat.StatsToConsole(args[0])
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
