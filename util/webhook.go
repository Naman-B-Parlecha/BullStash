package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	ColorGreen  = 0x00FF00 // Success
	ColorRed    = 0xFF0000 // Error
	ColorBlue   = 0x0000FF // Info
	ColorYellow = 0xFFFF00 // Warning
	ColorPurple = 0x800080 // Special
)

type DiscordEmbed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"`
}

type DiscordMessage struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

func SendToDiscord(webhook string, message *DiscordMessage) error {
	jsonData, _ := json.Marshal(message)

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook error: %s", resp.Status)
	}
	return nil
}

func CallWebHook(text string, isError bool) {
	if err := godotenv.Load(); err != nil {
		WarningColor.Println("No .env file found")
	}
	webhookURL, exists := os.LookupEnv("DISCORD_WEBHOOK_URL")

	if !exists {
		WarningColor.Println("discord webhook not set")
	}

	if webhookURL == "" {
		WarningColor.Println("Discord webhook URL is empty. Skipping webhook call.")
		return
	}

	color := ColorGreen
	title := "Success"

	if isError {
		color = ColorRed
		title = "Error"
	}

	message := &DiscordMessage{
		Content: fmt.Sprintf("Task executed at %s", time.Now().Format(time.RFC1123)),
		Embeds: []DiscordEmbed{
			{
				Title:       title,
				Description: fmt.Sprintf("```\n%s\n```", text),
				Color:       color,
			},
		},
	}

	if err := SendToDiscord(webhookURL, message); err != nil {
		ErrorColor.Printf("Failed to send to Discord: %v\n", err)
	}
}
