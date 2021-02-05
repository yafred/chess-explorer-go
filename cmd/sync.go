package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yafred/chess-explorer/internal/sync"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Download recent games for all users in database",
	Long:  `Download recent games for all users in database`,
	Run: func(cmd *cobra.Command, args []string) {
		sync.All()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
