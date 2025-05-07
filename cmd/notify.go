/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Add your webhook for alerting",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		discordWebhook, _ := cmd.Flags().GetString("discord")
		if discordWebhook == "" {
			fmt.Println("Please provide a Discord webhook URL.")
			return
		}

		filecontent, err := os.ReadFile(".env")
		if err != nil {
			fmt.Println("Error reading .env file:", err)
			return
		}
		content := string(filecontent)
		if len(content) > 0 && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += "\nDISCORD_WEBHOOK_URL=" + discordWebhook + "\n"

		err = os.WriteFile(".env", []byte(content), 0644)
		if err != nil {
			fmt.Println("Error writing to .env file:", err)
			return
		}

		fmt.Printf("Notification will be sent to Discord webhook: %s\n", discordWebhook)
		fmt.Println("WebHook added successfully!")
	},
}

func init() {
	rootCmd.AddCommand(notifyCmd)
	notifyCmd.Flags().String("discord", "", "Discord webhook URL")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
