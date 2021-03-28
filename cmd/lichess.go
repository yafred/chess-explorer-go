package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yafred/chess-explorer/internal/lichess"
)

var userToken string
var lichessPgn string

var lichessCmd = &cobra.Command{
	Use:   "lichess [user]",
	Short: "Download games for a given user from Lichess.org",
	Long:  `Download games for a given user from Lichess.org`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			lichess.DownloadGames(arg, lichessPgn)
		}
	},
}

func init() {
	rootCmd.AddCommand(lichessCmd)

	lichessCmd.Flags().StringVar(&userToken, "token", "", "your lichess.org personal API access token")
	lichessCmd.Flags().StringVar(&lichessPgn, "keep", "", "file where the PGN will be kept")

	// To be able to support the config file, we need to bind with viper (and read with viper.GetString())
	viper.BindPFlag("lichess-token", lichessCmd.Flags().Lookup("token"))
}
