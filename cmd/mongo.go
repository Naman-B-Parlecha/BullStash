/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// mongoCmd represents the mongo command
var mongoCmd = &cobra.Command{
	Use:   "mongo",
	Short: "",
	Long:  `Add your monogdb credential to env for easy backups`,
	Run: func(cmd *cobra.Command, args []string) {
		mongoURI, _ := cmd.Flags().GetString("mongo-uri")
		dbname, _ := cmd.Flags().GetString("dbname")

		if mongoURI == "" || dbname == "" {
			util.ErrorColor.Printf("Please provide a valid MongoDB URI and database name.")
			util.CallWebHook("Please provide a valid MongoDB URI and database name.", true)
			return
		}

		util.InfoColor.Printf("Kindly Confirm the following details :"+"\n"+" mongo-uri = "+mongoURI+"\n", "dbname = "+dbname, "\n\nEnter y/n to confirm")

		var confirm string
		fmt.Scanln(&confirm)

		if confirm == "n" {
			util.CallWebHook("User cancelled the database setup", true)
			return
		}
		existingContent := ""
		if fileInfo, err := os.Stat(".env"); err == nil && fileInfo.Size() > 0 {
			contentBytes, err := os.ReadFile(".env")
			if err == nil {
				existingContent = string(contentBytes)
			}
		}

		file, err := os.Create(".env")
		if err != nil {
			util.CallWebHook("Failed to save your database details, try again", true)
			util.ErrorColor.Printf("Failed to save your database details, try again")
			return
		}

		content := existingContent
		if len(existingContent) > 0 && !strings.HasSuffix(existingContent, "\n") {
			content += "\n\n"
		}

		content += "MONGO_URI=" + mongoURI + "\n"
		content += "MONGO_DB_NAME=" + dbname + "\n"
		if _, err := file.WriteString(content); err != nil {
			util.CallWebHook("Failed to write database details to file, try again", true)
			util.ErrorColor.Printf("Failed to write database details to file, try again")
			return
		}
		defer file.Close()
		util.SuccessColor.Printf("Database details saved successfully to .env file.")
		util.CallWebHook("Database details saved successfully", false)
	},
}

func init() {
	rootCmd.AddCommand(mongoCmd)

	mongoCmd.Flags().String("mongo-uri", "", "MongoDB URI")
	mongoCmd.Flags().String("dbname", "", "Database name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mongoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mongoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
