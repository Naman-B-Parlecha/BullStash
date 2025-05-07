package postgres

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func Backup(output, dbname, host, user, password string, port int, compress bool, storage string) error {

	projectDir, err := os.Getwd()

	if err != nil {
		util.CallWebHook("Error getting current directory: "+err.Error(), true)
		fmt.Printf("Error getting current directory: %v\n", err)
		return err
	}
	folderName := filepath.Join(projectDir, output)
	if err := os.MkdirAll(folderName, 0755); err != nil {
		util.CallWebHook("Error creating output directory: "+err.Error(), true)
		fmt.Printf("Failed to create output directory: %v\n", err)
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s/%s_backup_%s.sql", output, dbname, timestamp)
	gzFileName := fileName + ".gz"

	sqlFile, err := os.Create(fileName)
	if err != nil {
		util.CallWebHook("Error creating backup file: "+err.Error(), true)
		fmt.Printf("Failed to create backup file: %v\n", err)
		return err
	}
	defer sqlFile.Close()

	dumpCmd := exec.Command("pg_dump",
		"-h", host,
		"-p", strconv.Itoa(port),
		"-U", user,
		"-d", dbname)
	dumpCmd.Env = append(os.Environ(), "PGPASSWORD="+password)
	dumpCmd.Stdout = sqlFile

	if err := dumpCmd.Run(); err != nil {
		util.CallWebHook("pg_dump failed: "+err.Error(), true)
		fmt.Printf("pg_dump failed: %v\n", err)
		os.Remove(fileName)
		return err
	}

	if compress {
		if err := util.CompressFile(fileName, gzFileName); err != nil {
			fmt.Printf("Compression failed: %v\n", err)
			return err
		}
		if err := os.Remove(fileName); err != nil {
			util.CallWebHook("Error removing uncompressed file: "+err.Error(), true)
			fmt.Printf("Warning: could not remove uncompressed file: %v\n", err)
		}

		util.CallWebHook("Backup created successfully at: "+gzFileName, false)
		fmt.Printf("Backup successfully created at: %s\n", gzFileName)
		client := resty.New()
		fileInfo, err := os.Stat(gzFileName)
		if err != nil {
			fmt.Printf("Error getting file size: %v\n", err)
		}
		fileSize := fileInfo.Size()

		_, err = client.R().SetBody(struct {
			DBType     string  `json:"dbtype"`
			BackupType string  `json:"backup_type"`
			Storage    string  `json:"storage"`
			Size       float64 `json:"size"`
		}{
			DBType:     "postgres",
			BackupType: "full",
			Storage:    "local",
			Size:       float64(fileSize),
		}).Post("http://localhost:8085/backups/size")

		fmt.Println("File size sent to monitoring service:", fileSize)
		if err != nil {
			fmt.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
		}
		return nil
	}

	// read the file and put it to the bucket if it cloud storage

	if storage == "cloud" {

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
			fmt.Printf("Error creating AWS session: %v\n", err)
			return err
		}
		s3Client := s3.New(sess)

		file, err := os.Open(fileName)
		if err != nil {
			util.CallWebHook("Error opening file: "+err.Error(), true)
			fmt.Printf("Error opening file: %v\n", err)
		}

		defer file.Close()
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(awsBucketName),
			Key:    aws.String(fileName),
		})

		if err != nil {
			util.CallWebHook("Error uploading file to S3: "+err.Error(), true)
			fmt.Printf("Error uploading file to S3: %v\n", err)
			return err
		}
	}

	client := resty.New()
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Printf("Error getting file size: %v\n", err)
	}
	fileSize := fileInfo.Size()

	_, err = client.R().SetBody(struct {
		DBType     string  `json:"dbtype"`
		BackupType string  `json:"backup_type"`
		Storage    string  `json:"storage"`
		Size       float64 `json:"size"`
	}{
		DBType:     "postgres",
		BackupType: "full",
		Storage:    "local",
		Size:       float64(fileSize),
	}).Post("http://localhost:8085/backups/size")

	fmt.Println("File size sent to monitoring service:", fileSize)
	if err != nil {
		fmt.Println("Error sending request:", err)
		util.CallWebHook("Error sending request: "+err.Error(), true)
	}

	util.CallWebHook("Backup created successfully at: "+fileName, false)
	fmt.Printf("Backup successfully created at: %s\n", fileName)

	if storage == "cloud" {
		if err := os.Remove(fileName); err != nil {
			util.CallWebHook("Error removing local file: "+err.Error(), true)
			fmt.Printf("Warning: could not remove local file: %v\n", err)
		}
	}

	return nil
}
