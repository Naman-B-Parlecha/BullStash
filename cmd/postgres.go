/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// postgresCmd represents the postgres command
var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		password, _ := cmd.Flags().GetString("password")
		user, _ := cmd.Flags().GetString("user")
		name, _ := cmd.Flags().GetString("dbname")
		port, _ := cmd.Flags().GetInt("port")

		fmt.Println("Kindly Confirm the following details :"+"\n"+"Host: "+host+"\n"+"Password: "+password+"\n"+"User: "+user+"\n"+"Name: "+name+"\n"+"Port: ", port, "\n\nEnter y/n to confirm")

		var confirm string
		fmt.Scanln(&confirm)

		if confirm == "n" {
			return
		}

		file, err := os.Create(".env")

		if err != nil {
			fmt.Println("Failed to save your database details, try again")
		}

		file.WriteString("POSTGRES_DB_HOST=" + host + "\n" + "POSTGRES_DB_PORT=" + strconv.Itoa(port) + "\n" + "POSTGRES_DB_USER=" + user + "\n" + "POSTGRES_DB_PASSWORD=" + password + "\n" + "POSTGRES_DB_NAME=" + name)
		file.Close()
	},
}

func init() {
	rootCmd.AddCommand(postgresCmd)

	postgresCmd.Flags().String("host", "localhost", "Host of the database")
	postgresCmd.Flags().Int("port", 5432, "Port of the database")
	postgresCmd.Flags().String("dbname", "", "Name of the database")
	postgresCmd.Flags().String("user", "postgres", "User of the database")
	postgresCmd.Flags().String("password", "postgres", "Password of the database")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postgresCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postgresCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
