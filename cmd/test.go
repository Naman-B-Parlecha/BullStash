/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/mongo"
	"github.com/Naman-B-Parlecha/BullStash/mysql"
	"github.com/Naman-B-Parlecha/BullStash/postgres"
	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "",
	Long:  `Test your database connectivity using this command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("checking connection to database")
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		mongoURI, _ := cmd.Flags().GetString("mongo_uri")

		start := time.Now()
		client := resty.New()

		if dbtype == "" {
			util.CallWebHook("Database type not specified", true)
			fmt.Println("Database type not specified")

			_, err := client.R().SetBody(struct {
				DBType string `json:"dbtype"`
			}{
				DBType: "Unknown",
			}).Post("http://localhost:8085/connections/failure")

			if err != nil {
				fmt.Println("Error sending request:", err)
				util.CallWebHook("Error sending request: "+err.Error(), true)
			}
			return
		}
		if dbtype == "postgres" {
			err := postgres.TestConnection(port, dbname, host, user, password)
			if err != nil {
				util.CallWebHook("Error connecting to database: "+err.Error(), true)
				fmt.Printf("Error connecting to database: %v\n", err)
				_, err := client.R().SetBody(struct {
					DBType string `json:"dbtype"`
				}{
					DBType: dbtype,
				}).Post("http://localhost:8085/connections/failure")

				if err != nil {
					fmt.Println("Error sending request:", err)
					util.CallWebHook("Error sending request: "+err.Error(), true)
				}
				return
			} else {
				util.CallWebHook("Database connection successful", false)
				fmt.Println("Database connection successful")
			}
		}
		if dbtype == "mysql" {
			err := mysql.TestConnection(dbname, user, password)
			if err != nil {
				util.CallWebHook("Error connecting to database: "+err.Error(), true)
				fmt.Printf("Error connecting to database: %v\n", err)
				_, err := client.R().SetBody(struct {
					DBType string `json:"dbtype"`
				}{
					DBType: dbtype,
				}).Post("http://localhost:8085/connections/failure")

				if err != nil {
					fmt.Println("Error sending request:", err)
					util.CallWebHook("Error sending request: "+err.Error(), true)
				}
				return
			} else {
				util.CallWebHook("Database connection successful", false)
				fmt.Println("Database connection successful")
			}
		}

		if dbtype == "mongo" {
			err := mongo.TestConnection(mongoURI)
			if err != nil {
				util.CallWebHook("Error connecting to database: "+err.Error(), true)
				fmt.Printf("Error connecting to database: %v\n", err)
				_, err := client.R().SetBody(struct {
					DBType string `json:"dbtype"`
				}{
					DBType: dbtype,
				}).Post("http://localhost:8085/connections/failure")

				if err != nil {
					fmt.Println("Error sending request:", err)
					util.CallWebHook("Error sending request: "+err.Error(), true)
				}
				return
			} else {
				util.CallWebHook("Database connection successful", false)
				fmt.Println("Database connection successful")
			}
		}

		_, err := client.R().SetBody(struct {
			DBType string `json:"dbtype"`
		}{
			DBType: dbtype,
		}).Post("http://localhost:8085/connections/success")

		if err != nil {
			fmt.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
		}

		end := time.Since(start)

		_, err = client.R().SetBody(struct {
			DBType   string  `json:"dbtype"`
			Duration float64 `json:"duration"`
		}{
			DBType:   dbtype,
			Duration: float64(end.Milliseconds()),
		}).Post("http://localhost:8085/connections/latency")

		if err != nil {
			fmt.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
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
	testCmd.Flags().String("mongo_uri", "", "MongoDB URI")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
