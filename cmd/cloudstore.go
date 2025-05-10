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

// cloudstoreCmd represents the cloudstore command
var cloudstoreCmd = &cobra.Command{
	Use:   "cloudstore",
	Short: "A command to manage cloud storage",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Here you will define your command's action.
		region, _ := cmd.Flags().GetString("region")
		bucket, _ := cmd.Flags().GetString("bucket")
		accessKey, _ := cmd.Flags().GetString("access-key")
		secretKey, _ := cmd.Flags().GetString("secret-key")

		if region == "" || bucket == "" || accessKey == "" || secretKey == "" {
			util.WarningColor.Println("All flags (region, bucket, access-key, secret-key) are required.")
			return
		}
		util.SuccessColor.Printf("Region: %s\n", region)
		util.SuccessColor.Printf("Bucket: %s\n", bucket)
		util.SuccessColor.Printf("Access Key: %s\n", accessKey)
		util.SuccessColor.Printf("Secret Key: %s\n", secretKey)

		fileContent, err := os.ReadFile(".env")
		if err != nil {
			util.ErrorColor.Println("Error reading .env file:", err)
			return
		}
		content := string(fileContent)
		if len(content) > 0 && !strings.HasSuffix(content, "\n") {
			content += "\n\n"
		}

		content += fmt.Sprintf("CLOUD_REGION=%s\n", region)
		content += fmt.Sprintf("CLOUD_BUCKET=%s\n", bucket)
		content += fmt.Sprintf("CLOUD_ACCESS_KEY=%s\n", accessKey)
		content += fmt.Sprintf("CLOUD_SECRET_KEY=%s\n", secretKey)

		err = os.WriteFile(".env", []byte(content), 0644)
		if err != nil {
			util.SuccessColor.Println("Error writing to .env file:", err)
			return
		}

		util.SuccessColor.Printf("Cloud storage configuration added to .env file:\n")
		util.CallWebHook("Cloud storage configuration added to .env file:\n", false)
	},
}

func init() {
	rootCmd.AddCommand(cloudstoreCmd)
	cloudstoreCmd.Flags().StringP("region", "r", "", "Specify the region for the cloud store")
	cloudstoreCmd.Flags().StringP("bucket", "b", "", "Specify the bucket name for the cloud store")
	cloudstoreCmd.Flags().StringP("access-key", "a", "", "Specify the access key for the cloud store")
	cloudstoreCmd.Flags().StringP("secret-key", "s", "", "Specify the secret key for the cloud store")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloudstoreCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloudstoreCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
