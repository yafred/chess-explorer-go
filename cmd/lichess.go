package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/lichess"
)

var userToken string
var timeout int

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

	lichessCmd.Flags().StringVar(&userToken, "token", "", "your lichess.org personal API access token")
	lichessCmd.Flags().IntVar(&timeout, "timeout", 300, "timeout value in seconds when downloading games from lichess.org")

	// To be able to support the config file, we need to bind with viper (and read with viper.GetString())
	viper.BindPFlag("lichess-token", lichessCmd.Flags().Lookup("token"))
	viper.BindPFlag("lichess-timeout", lichessCmd.Flags().Lookup("timeout"))
}
