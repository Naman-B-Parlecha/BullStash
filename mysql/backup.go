package mysql

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Naman-B-Parlecha/BullStash/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func Backup(output, dbname, user, password string, compress bool, storage string) error {

	projectDir, err := os.Getwd()

	if err != nil {
		return err
	}
	folderName := filepath.Join(projectDir, output)
	if err := os.MkdirAll(folderName, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s/%s_backup_%s.sql", output, dbname, timestamp)
	gzFileName := fileName + ".gz"

	sqlFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer sqlFile.Close()

	dumpCmd := exec.Command("mysqldump",
		"-u", user,
		"--password="+password,
		dbname)

	dumpCmd.Stdout = sqlFile

	if err := dumpCmd.Run(); err != nil {
		os.Remove(fileName)
		return fmt.Errorf("mysqldump failed: %v", err)
	}

	if compress {
		if err := util.CompressFile(fileName, gzFileName); err != nil {
			return err
		}
		if err := os.Remove(fileName); err != nil {
			util.CallWebHook("Error removing uncompressed file: "+err.Error(), true)
			util.ErrorColor.Printf("Warning: could not remove uncompressed file: %v\n", err)
		}

		util.CallWebHook("Backup created successfully at: "+gzFileName, false)
		util.SuccessColor.Printf("Backup successfully created at: %s\n", gzFileName)

		client := resty.New()
		fileInfo, err := os.Stat(gzFileName)
		if err != nil {
			util.ErrorColor.Printf("Error getting file size: %v\n", err)
		}
		fileSize := fileInfo.Size()

		_, err = client.R().SetBody(struct {
			DBType     string  `json:"dbtype"`
			BackupType string  `json:"backup_type"`
			Storage    string  `json:"storage"`
			Size       float64 `json:"size"`
		}{
			DBType:     "mysql",
			BackupType: "full",
			Storage:    "local",
			Size:       float64(fileSize),
		}).Post("http://localhost:8085/backups/size")

		util.SuccessColor.Println("File size sent to monitoring service:", fileSize)
		if err != nil {
			util.ErrorColor.Println("Error sending request:", err)
			util.CallWebHook("Error sending request: "+err.Error(), true)
		}

		return nil
	}

	if storage == "cloud" {
		godotenv.Load(".env")

		cloudRegion := os.Getenv("CLOUD_REGION")
		cloudAccessKey := os.Getenv("CLOUD_ACCESS_KEY")
		cloudSecretKey := os.Getenv("CLOUD_SECRET_KEY")
		cloudBucketName := os.Getenv("CLOUD_BUCKET")

		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(cloudRegion),
			Credentials: credentials.NewStaticCredentials(
				cloudAccessKey, cloudSecretKey, ""),
		})

		if err != nil {
			util.CallWebHook("Error creating AWS session: "+err.Error(), true)
			util.ErrorColor.Printf("Error creating AWS session: %v\n", err)
		}

		s3Client := s3.New(sess)
		file, err := os.Open(fileName)
		if err != nil {
			util.CallWebHook("Error opening file: "+err.Error(), true)
			util.ErrorColor.Printf("Error opening file: %v\n", err)
		}

		defer file.Close()
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(cloudBucketName),
			Key:    aws.String(fileName),
		})

		if err != nil {
			util.CallWebHook("Error uploading file to S3: "+err.Error(), true)
			util.ErrorColor.Printf("Error uploading file to S3: %v\n", err)
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
		DBType:     "mysql",
		BackupType: "full",
		Storage:    "local",
		Size:       float64(fileSize),
	}).Post("http://localhost:8085/backups/size")

	util.SuccessColor.Println("File size sent to monitoring service:", fileSize)
	if err != nil {
		fmt.Println("Error sending request:", err)
		util.CallWebHook("Error sending request: "+err.Error(), true)
	}

	util.CallWebHook("Backup created successfully at: "+fileName, false)
	util.ErrorColor.Printf("Backup successfully created at: %s\n", fileName)

	if storage == "cloud" {
		if err := os.Remove(fileName); err != nil {
			util.CallWebHook("Error removing local file: "+err.Error(), true)
			util.WarningColor.Printf("Warning: could not remove local file: %v\n", err)
		}
	}
	return nil
}
