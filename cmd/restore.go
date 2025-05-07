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

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "A brief description of your command",
	Long:  `Restore your database using this command`,
	Run: func(cmd *cobra.Command, args []string) {
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		input, _ := cmd.Flags().GetString("input")
		mongo_uri, _ := cmd.Flags().GetString("mongo-uri")
		isDrop, _ := cmd.Flags().GetBool("drop")
		iscompressed, _ := cmd.Flags().GetBool("isCompressed")

		start := time.Now()

		client := resty.New()
		if dbtype == "" {
			util.CallWebHook("Please enter a valid database type", true)
			fmt.Println("Enter a valid Database Type")
			_, err := client.R().SetBody(struct {
				DBType  string `json:"dbtype"`
				Storage string `json:"storage"`
			}{
				DBType:  "Unknown",
				Storage: "local",
			}).Post("http://localhost:8085/restore/failure")

			if err != nil {
				fmt.Println("Error sending request:", err)
				util.CallWebHook("Error sending request: "+err.Error(), true)
			}
			return
		}

		if dbtype == "postgres" {
			err := postgres.Restore(dbname, input, host, user, password, port)
			if err != nil {
				util.CallWebHook("Error restoring database: "+err.Error(), true)
				fmt.Printf("Error restoring database: %v\n", err)

				_, err := client.R().SetBody(struct {
					DBType  string `json:"dbtype"`
					Storage string `json:"storage"`
				}{
					DBType:  dbtype,
					Storage: "local",
				}).Post("http://localhost:8085/restore/failure")

				if err != nil {
					fmt.Println("Error sending request:", err)
					util.CallWebHook("Error sending request: "+err.Error(), true)
				}
				return
			}
		}

		if dbtype == "mysql" {
			err := mysql.Restore(dbname, input, user, password)
			if err != nil {
				util.CallWebHook("Error restoring database: "+err.Error(), true)
				fmt.Printf("Error restoring database: %v\n", err)

				_, err := client.R().SetBody(struct {
					DBType  string `json:"dbtype"`
					Storage string `json:"storage"`
				}{
					DBType:  dbtype,
					Storage: "local",
				}).Post("http://localhost:8085/restore/failure")

				if err != nil {
					fmt.Println("Error sending request:", err)
					util.CallWebHook("Error sending request: "+err.Error(), true)
				}
				return
			}
		}

		if dbtype == "mongo" {
			err := mongo.Restore(mongo_uri, input, isDrop, iscompressed)
			if err != nil {
				util.CallWebHook("Error restoring database: "+err.Error(), true)
				fmt.Printf("Error restoring database: %v\n", err)
				_, err := client.R().SetBody(struct {
					DBType  string `json:"dbtype"`
					Storage string `json:"storage"`
				}{
					DBType:  dbtype,
					Storage: "local",
				}).Post("http://localhost:8085/restore/failure")

				if err != nil {
					fmt.Println("Error sending request:", err)
					util.CallWebHook("Error sending request: "+err.Error(), true)
				}
				return
			}
		}

		_, err := client.R().SetBody(struct {
			DBType  string `json:"dbtype"`
			Storage string `json:"storage"`
		}{
			DBType:  dbtype,
			Storage: "local",
		}).Post("http://localhost:8085/restore/success")

		if err != nil {
			fmt.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
		}

		end := time.Since(start)
		fmt.Println("Time taken to restore the database:", end.Milliseconds())
		_, err = client.R().SetBody(struct {
			DBType   string  `json:"dbtype"`
			Storage  string  `json:"storage"`
			Duration float64 `json:"duration"`
		}{
			DBType:   dbtype,
			Storage:  "local",
			Duration: float64(end.Milliseconds()),
		}).Post("http://localhost:8085/restore/duration")

		if err != nil {
			fmt.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
		}

	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().String("dbtype", "postgres", "Database type")
	restoreCmd.Flags().String("host", "localhost", "Database host")
	restoreCmd.Flags().Int("port", 5432, "Database port")
	restoreCmd.Flags().String("user", "postgres", "Database user")
	restoreCmd.Flags().String("password", "", "Database password")
	restoreCmd.Flags().String("dbname", "", "Database name")
	restoreCmd.Flags().String("input", "", "Input file")
	restoreCmd.Flags().Bool("drop", true, "do u want to drop ur mongo collections before restore??")
	restoreCmd.Flags().String("mongo-uri", "", "MongoDB URI that u want to restore to")
	restoreCmd.Flags().Bool("isCompressed", false, "Is your dump files compressed")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
