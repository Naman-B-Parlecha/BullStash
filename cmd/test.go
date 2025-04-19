/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Naman-B-Parlecha/BullStash/postgres"
	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("checking connection to database")
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")

		if dbtype == "postgres" {
			err := postgres.TestConnection(port, dbname, host, user, password)
			if err != nil {
				util.CallWebHook("Error connecting to database: "+err.Error(), true)
				fmt.Printf("Error connecting to database: %v\n", err)
			} else {
				util.CallWebHook("Database connection successful", false)
				fmt.Println("Database connection successful")
			}
		}
		if dbtype == "mysql" {
			// will implement mysql test connection logic here

		}

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().String("dbtype", "postgres", "Database type")
	testCmd.Flags().String("host", "localhost", "Database host")
	testCmd.Flags().Int("port", 5432, "Database port")
	testCmd.Flags().String("user", "postgres", "Database user")
	testCmd.Flags().String("password", "", "Database password")
	testCmd.Flags().String("dbname", "postgres", "Database name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
