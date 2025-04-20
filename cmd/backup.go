/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/Naman-B-Parlecha/BullStash/internal/config"
	"github.com/Naman-B-Parlecha/BullStash/mongo"
	"github.com/Naman-B-Parlecha/BullStash/mysql"
	"github.com/Naman-B-Parlecha/BullStash/postgres"
	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dbtype, _ := cmd.Flags().GetString("dbtype")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		dbname, _ := cmd.Flags().GetString("dbname")
		output, _ := cmd.Flags().GetString("output")
		compress, _ := cmd.Flags().GetBool("compress")
		isCron, _ := cmd.Flags().GetBool("isCron")
		mongoURI, _ := cmd.Flags().GetString("mongo-uri")
		if dbtype == "" {
			util.CallWebHook("Please enter a valid database type", true)
			fmt.Println("Enter a valid Database Type")
			return
		}

		if dbname == "postgres" {
			if isCron {
				postgresConfig := config.GetPostgresConfig()
				host = postgresConfig.HOST
				port, _ = strconv.Atoi(postgresConfig.PORT)
				user = postgresConfig.USER
				password = postgresConfig.PASSWORD
				dbname = postgresConfig.DBNAME
			}

			err := postgres.Backup(output, dbname, host, user, password, port, compress)
			if err != nil {
				util.CallWebHook("Error creating backup: "+err.Error(), true)
				fmt.Printf("Error creating backup: %v\n", err)
			}
		}

		if dbtype == "mysql" {
			fmt.Println("Creating backup for MySQL database")
			if isCron {
				mysqlConfig := config.GetMySqlConfig()
				user = mysqlConfig.USER
				password = mysqlConfig.PASSWORD
				dbname = mysqlConfig.DBNAME
			}

			err := mysql.Backup(output, dbname, user, password, compress)

			if err != nil {
				util.CallWebHook("Error creating backup: "+err.Error(), true)
				fmt.Printf("Error creating backup: %v\n", err)
			}
		}

		if dbtype == "mongo" {
			if isCron {
				mongoConfig := config.GetMongoConfig()
				mongoURI = mongoConfig.MONGO_URI
				dbname = mongoConfig.DBNAME
			}
			err := mongo.Backup(mongoURI, dbname, output, compress)
			if err != nil {
				util.CallWebHook("Error creating backup: "+err.Error(), true)
				fmt.Printf("Error creating backup: %v\n", err)
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().String("dbtype", "postgres", "Type of database")
	backupCmd.Flags().String("host", "localhost", "Host of the database")
	backupCmd.Flags().Int("port", 0, "Port of the database")
	backupCmd.Flags().String("user", "", "User of the database")
	backupCmd.Flags().String("password", "", "Password of the database")
	backupCmd.Flags().String("dbname", "postgres", "Name of the database")
	backupCmd.Flags().String("backup-type", "full", "Type of backup such as full, incremental, differential")
	backupCmd.Flags().String("output", "backups", "Path you want to store the backup")
	backupCmd.Flags().Bool("compress", false, "Compress the backup file")
	backupCmd.Flags().String("storage", "local", "Storage type such as local, s3, gcs, goodle_drive")
	backupCmd.Flags().String("cloud-bucket", "", "Cloud bucket name")
	backupCmd.Flags().String("cloud-region", "asia-pacific-1", "Cloud region")
	backupCmd.Flags().Bool("isCron", false, "Is this a cron job")

	backupCmd.Flags().String("mongo-uri", "", "MongoDB URI")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
