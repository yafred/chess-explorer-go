package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	stat "github.com/yafred/chess-com/internal/chess-stat"
)

var csvFilePath string
var pgnFilePath string
var cachePath string
var cacheRefresh bool
var showStats bool

//var cacheRefresh bool

var playerCmd = &cobra.Command{
	Use:   "player [username]",
	Short: "Creates stats for a chess.com player",
	Long:  `Creates stats for a chess.com player based on games downloaded from https://chess.com/`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if showStats || (csvFilePath == "" && pgnFilePath == "") {
			stat.StatsToConsole(strings.ToLower(args[0]), viper.GetString("cache-path"), cacheRefresh)
		}
		if csvFilePath != "" {
			stat.CreateCsvFile(strings.ToLower(args[0]), viper.GetString("cache-path"), cacheRefresh, csvFilePath)
		}
		if pgnFilePath != "" {
			stat.CreatePgnFile(strings.ToLower(args[0]), viper.GetString("cache-path"), cacheRefresh, pgnFilePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(playerCmd)

	playerCmd.Flags().BoolVarP(&showStats, "show-stats", "s", false, "Show statistics")
	playerCmd.Flags().StringVarP(&cachePath, "cache-path", "c", "", "Folder where downloaded data should be kept (data will not be kept if flag absent)")
	playerCmd.MarkFlagDirname("cache-path")
	playerCmd.Flags().BoolVarP(&cacheRefresh, "cache-refresh", "r", false, "Refresh cache before executing command (if flag absent, existing data will be used)")
	playerCmd.Flags().StringVarP(&csvFilePath, "csv", "", "", "Creates a csv file")
	playerCmd.MarkFlagFilename("csv")
	playerCmd.Flags().StringVarP(&pgnFilePath, "pgn", "", "", "Creates a pgn file")
	playerCmd.MarkFlagFilename("pgn")

	// To be able to support the config file, we need to bind with viper (and read with viper.GetString())
	viper.BindPFlag("cache-path", playerCmd.Flags().Lookup("cache-path"))
}
