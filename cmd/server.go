package cmd

import (
	"github.com/spf13/cobra"
	server "github.com/yafred/chess-stat/web/app"
)

var serverPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a web server on port 8080",
	Long:  `Web server and API server`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(serverPort)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntVar(&serverPort, "server-port", 8080, "Server http port (default is 8080)")
}
