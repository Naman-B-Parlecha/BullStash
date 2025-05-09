/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "BullStash",
	Short: "A CLI tool for automated managing of your database backups and restores.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
██████╗ ██╗   ██╗██╗     ██╗      ███████╗████████╗ █████╗ ███████╗██╗  ██╗
██╔══██╗██║   ██║██║     ██║      ██╔════╝╚══██╔══╝██╔══██╗██╔════╝██║  ██║
██████╔╝██║   ██║██║     ██║      ███████╗   ██║   ███████║███████╗███████║
██╔══██╗██║   ██║██║     ██║	  ╚════██║   ██║   ██╔══██║╚════██║██╔══██║
██████╔╝╚██████╔╝███████╗███████╗ ███████║   ██║   ██║  ██║███████║██║  ██║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝ ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
		`)
		util.InfoColor.Println("Welcome to BullStash! A CLI tool for managing your database backups and restores.")
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.BullStash.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
