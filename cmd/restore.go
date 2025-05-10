/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/mongo"
	"github.com/Naman-B-Parlecha/BullStash/mysql"
	"github.com/Naman-B-Parlecha/BullStash/postgres"
	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
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
		storage, _ := cmd.Flags().GetString("storage")

		if storage == "cloud" {
			// first i will list all the items so that i can get the file name
			// then i will download the file to a temp location
			// then i will restore the file

			godotenv.Load(".env")
			awsAccessKey := os.Getenv("CLOUD_ACCESS_KEY")
			awsSecretKey := os.Getenv("CLOUD_SECRET_KEY")
			awsBucketName := os.Getenv("CLOUD_BUCKET")
			awsBucketRegion := os.Getenv("CLOUD_REGION")

			sess, err := session.NewSession(&aws.Config{
				Region:      aws.String(awsBucketRegion),
				Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
			})

			if err != nil {
				util.CallWebHook("Error creating AWS session: "+err.Error(), true)
				fmt.Println("Error creating AWS session:", err)
				return
			}

			s3Client := s3.New(sess)

			resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket: aws.String(awsBucketName),
			})

			if err != nil {
				util.CallWebHook("Error listing S3 objects: "+err.Error(), true)
				fmt.Println("Error listing S3 objects:", err)
				return
			}

			for _, item := range resp.Contents {
				fmt.Println("File Name: ", *item.Key)
			}

			downloader := s3manager.NewDownloader(sess)

			fmt.Printf("Enter the file name to download: ")
			var fileName string
			fmt.Scanln(&fileName)

			if fileName == "" {
				util.CallWebHook("Please enter a valid file name", true)
				fmt.Println("Enter a valid file name")
				return
			}

			filePath := fileName

			if strings.Contains(fileName, "/") {
				dirPath := filepath.Dir(fileName)
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						util.CallWebHook("Error creating directories: "+err.Error(), true)
						fmt.Println("Error creating directories:", err)
						return
					}
				}
			}
			file, err := os.Create(filePath)
			if err != nil {
				util.CallWebHook("Error creating file: "+err.Error(), true)
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()
			_, err = downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(awsBucketName),
					Key:    aws.String(fileName),
				})
			if err != nil {
				util.CallWebHook("Error downloading file: "+err.Error(), true)
				fmt.Println("Error downloading file:", err)
				return
			}

			fmt.Printf("Successfully downloaded %s to %s\n", fileName, filePath)
			util.CallWebHook(fmt.Sprintf("Successfully downloaded %s to %s", fileName, filePath), false)

			input = filePath
		}
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
	restoreCmd.Flags().String("storage", "local", "Storage type")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// restoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// restoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
