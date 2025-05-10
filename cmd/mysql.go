/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Add your mysql credential to env for easy backups",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		password, _ := cmd.Flags().GetString("password")
		user, _ := cmd.Flags().GetString("user")
		name, _ := cmd.Flags().GetString("dbname")
		port, _ := cmd.Flags().GetInt("port")

		if name == "" {
			util.CallWebHook("Please enter a valid database name", true)
			util.WarningColor.Println("Enter a valid Dabtabase Name")
			return
		}
		util.InfoColor.Println("Kindly Confirm the following details :"+"\n"+"Host: "+host+"\n"+"Password: "+password+"\n"+"User: "+user+"\n"+"Name: "+name+"\n"+"Port: ", port, "\n\nEnter y/n to confirm")

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
			util.ErrorColor.Println("Failed to save your database details, try again")
			return
		}

		content := existingContent
		if len(existingContent) > 0 && !strings.HasSuffix(existingContent, "\n") {
			content += "\n\n"
		}
		content += "MYSQL_DB_HOST=" + host + "\n"
		content += "MYSQL_DB_USER=" + user + "\n"
		content += "MYSQL_DB_PASSWORD=" + password + "\n"
		content += "MYSQL_DB_NAME=" + name + "\n"
		content += "MYSQL_DB_PORT=" + strconv.Itoa(port) + "\n"

		if _, err := file.WriteString(content); err != nil {
			util.CallWebHook("Failed to write database details to file, try again", true)
			util.ErrorColor.Println("Failed to write database details to file, try again")
			return
		}
		defer file.Close()
		util.CallWebHook("Database details saved successfully", false)
		util.SuccessColor.Println("Database details saved successfully to .env file.")
	},
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.Flags().StringP("host", "H", "localhost", "Host of the database")
	mysqlCmd.Flags().StringP("password", "p", "", "Password of the database")
	mysqlCmd.Flags().StringP("user", "u", "root", "User of the database")
	mysqlCmd.Flags().StringP("dbname", "d", "", "Name of the database")
	mysqlCmd.Flags().IntP("port", "P", 3306, "Port of the database")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mysqlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mysqlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
